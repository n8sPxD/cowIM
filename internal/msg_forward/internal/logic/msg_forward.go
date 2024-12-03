package logic

import (
	"context"
	"errors"
	"github.com/n8sPxD/cowIM/internal/msg_forward/gossip"
	"github.com/segmentio/kafka-go/compress"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"

	"github.com/n8sPxD/cowIM/internal/common/constant"
	"github.com/n8sPxD/cowIM/internal/common/message/front"
	"github.com/n8sPxD/cowIM/internal/common/message/inside"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/svc"
	"github.com/segmentio/kafka-go"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type MsgForwarder struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	Routes  *gossip.Server
	offline chan OfflineCacheMessage // 隔一段时间清一次，隔多久取决于Gossip集群大小和同步延迟

	MsgForwarder *kafka.Reader
	MsgSender    *map[int]*kafka.Writer
	MsgDBSaver   *kafka.Writer
}

// NewMsgSenderPool 创建消息发送Sender池
func NewMsgSenderPool(brokers []string, count int) *map[int]*kafka.Writer {
	// TODO: 当前只能固定数量，后续改进可以通过服务发现获取WebsocketServer数量，动态创建和销毁Sender
	pool := make(map[int]*kafka.Writer)
	for i := 1; i <= count; i++ {
		writer := &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        "websocket-server-" + strconv.Itoa(i),
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond, // 低超时时间
			RequiredAcks: kafka.RequireOne,      // 仅等待 Leader 确认
			Compression:  compress.Zstd,         // Zstd压缩
			Async:        true,                  // 启用异步写入
			MaxAttempts:  1,                     // 限制重试次数
		}
		pool[i] = writer
	}
	return &pool
}

func NewMsgForwarder(ctx context.Context, svcCtx *svc.ServiceContext) *MsgForwarder {
	return &MsgForwarder{
		ctx:     ctx,
		svcCtx:  svcCtx,
		Routes:  gossip.MustNewServer(svcCtx.Discov, svcCtx.Regist, svcCtx.Config.RPCPort, int(svcCtx.Config.WorkID), 3, 3, 3),
		offline: make(chan OfflineCacheMessage),
		MsgForwarder: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        svcCtx.Config.MsgForwarder.Brokers,
			Topic:          svcCtx.Config.MsgForwarder.Topic,
			GroupID:        "msg-fwd",
			StartOffset:    kafka.LastOffset,
			MinBytes:       1,                      // 最小拉取字节数
			MaxBytes:       10e3,                   // 最大拉取字节数（10KB）
			MaxWait:        100 * time.Millisecond, // 最大等待时间
			CommitInterval: 500 * time.Millisecond, // 提交间隔
		}),
		MsgSender: NewMsgSenderPool(svcCtx.Config.MsgSender.Brokers, 1), // 最后的count是WebsocketServer的数量
		MsgDBSaver: &kafka.Writer{
			Addr:         kafka.TCP(svcCtx.Config.MsgDBSaver.Brokers...),
			Topic:        svcCtx.Config.MsgDBSaver.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond, // 低超时时间
			RequiredAcks: kafka.RequireOne,      // 仅等待 Leader 确认
			Compression:  compress.Zstd,         // Zstd压缩
			Async:        true,                  // 启用异步写入
			MaxAttempts:  1,                     // 限制重试次数
		},
	}
}

func (l *MsgForwarder) Close() {
	l.MsgForwarder.Close()
}

func (l *MsgForwarder) Start() {
	// 初始化id生成器
	options := idgen.NewIdGeneratorOptions(l.svcCtx.Config.WorkID)
	idgen.SetIdGenerator(options)

	go l.Routes.Start()

	for {
		msg, err := l.MsgForwarder.ReadMessage(l.ctx) // 这里的msg是kafka.Message
		if err != nil {
			logx.Error("[Start] Reading message error: ", err)
			break
		}
		go l.Consume(msg.Value, time.Now())
	}
}

// Consume 接收从 Websocket Server的消息，处理后再进行转发
func (l *MsgForwarder) Consume(protobuf []byte, now time.Time) {
	// 传过来的消息是序列化过的，先反序列化
	var (
		msg   front.Message
		oldId string
	)
	err := proto.Unmarshal(protobuf, &msg)
	if err != nil {
		logx.Error("[Consume] Unmarshal message failed, error: ", err)
		return
	}

	// 隔离系统信息
	if msg.From != constant.USER_SYSTEM && msg.To != constant.USER_SYSTEM {
		// 检查消息重复性
		if ok, err := l.svcCtx.Redis.CheckDuplicateMessage(l.ctx, msg.Id); err != nil {
			logx.Error("[Consume] Check duplicate message from redis failed, error: ", err)
			return
		} else if ok { // 消息是重复的
			logx.Infof("[Consume] Message from %d with ID \"%s\" sent repeated", msg.From, msg.Id)
			return
		}
		// 保存以前的uuid,分配消息ID
		oldId, msg.Id = msg.Id, strconv.FormatInt(idgen.NextId(), 10)
		current := now.Unix()
		msg.Timestamp = current
		// 重新序列化
		protobuf, err = proto.Marshal(&msg)
		// 异步存库
		go l.sendRecordMsgToDB(&msg, now) // 漫游库
		go l.sendTimelineToDB(&msg, now)  // timeline
		// Ack消息
		go l.replyAckMessage(&msg, oldId)
	}

	// 进行基于消息类型的消息处理
	switch msg.Type {
	case constant.SINGLE_CHAT:
		go l.singleChat(&msg, protobuf)
	case constant.GROUP_CHAT:
		go l.groupChat(&msg, protobuf)
	case constant.BIG_GROUP_CHAT:
		go l.bigGroupChat(&msg)
	case constant.SYSTEM_INFO:
		go l.systemOperation(&msg, protobuf)
	default:
		logx.Error("[Consume] Wrong message type, Type is: ", msg.Type)
		return
	}
}

func (l *MsgForwarder) _packageMessageAndSend(protobuf []byte, id uint32, msgID string, msgType uint32, route int) {
	msg := inside.Message{
		To:       id,
		MsgId:    msgID,
		Protobuf: protobuf,
		Type:     msgType,
	}
	msgByte, err := proto.Marshal(&msg)
	if err != nil {
		logx.Error("[packageMessageAndSend] Marshal message failed, error: ", err)
		return
	}

	km := kafka.Message{Value: msgByte}
	err = (*l.MsgSender)[route].WriteMessages(l.ctx, km)
	if err != nil {
		logx.Error(
			"[packageMessage] Push message to Websocket-server MQ failed, error: ",
			err,
		)
	}
}

type OfflineCacheMessage struct {
	protobuf []byte
	id       uint32 // 用户ID
	msgID    string
	msgType  uint32
}

func (l *MsgForwarder) packageMessageAndSend(protobuf []byte, id uint32, msgID string, msgType uint32) {
	if router, ok := l.Routes.Node.Data[int32(id)]; ok {
		l._packageMessageAndSend(protobuf, id, msgID, msgType, int(router.Value))
	} else {
		// 没找到不代表没上线，防止消息误发的补偿机制
		/*
			message := OfflineCacheMessage{
				protobuf: protobuf,
				id:       id,
				msgID:    msgID,
				msgType:  msgType,
			}
			l.offline <- message
		*/
		/*
				 TODO: 完善消息缓存以及补偿机制，具体：定时消费l.offline，时间大概为Gossip同步周期+1秒，
			            为空时清空计时器，并且不计时，队列中有元素开始计时。MsgForwarder服务Start时启动异步消费任务
		*/
	}
}

// 单聊处理
func (l *MsgForwarder) singleChat(msg *front.Message, protobuf []byte) {
	l.packageMessageAndSend(protobuf, msg.To, msg.Id, msg.MsgType)
}

// 群聊处理，主动推
func (l *MsgForwarder) groupChat(msg *front.Message, protobuf []byte) {
	var ids any
	if members, err := l.svcCtx.Redis.GetGroupMembers(context.Background(), *msg.Group); err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存中还没有当前群聊的成员信息
			var sqlMembers []uint32
			if sqlMembers, err2 := l.svcCtx.MySQL.GetGroupMemberIDs(context.Background(), *msg.Group); err2 != nil {
				// 从SQL现查
				logx.Error("[groupChat] Get group members from mysql failed, error: ", err2)
				return
			} else if err2 := l.svcCtx.Redis.AddGroupMembers(context.Background(), *msg.Group, sqlMembers); err2 != nil {
				// MySQL查询没有问题，插缓存
				logx.Error("[groupChat] Add group members cache to redis failed, error: ", err2)
				return
			}
			// 查询没问题，redis没犯病，进入发送消息流程
			ids = sqlMembers
		} else {
			// redis错误
			logx.Error("[groupChat] Get group members from redis failed, error: ", err)
			return
		}
	} else {
		// redis有缓存，直接用
		ids = members
	}

	checkAndSend := func(current uint32, protobuf []byte, message *front.Message) {
		if current == msg.From {
			return
		} else {
			l.packageMessageAndSend(protobuf, msg.To, msg.Id, msg.MsgType)
		}
	}

	switch ids := ids.(type) {
	case []string:
		for _, member := range ids {
			realMember, _ := strconv.Atoi(member)
			checkAndSend(uint32(realMember), protobuf, msg)
		}
	case []int:
		for _, member := range ids {
			checkAndSend(uint32(member), protobuf, msg)
		}
	default:
		logx.Error("[sendMessageInGroup] Wrong members type!")
		return
	}
}

func (l *MsgForwarder) bigGroupChat(msg *front.Message) {
}

func (l *MsgForwarder) replyAckMessage(sender *front.Message, oldId string) {
	// 封装消息体
	reply := front.Message{
		Id:        oldId,
		From:      constant.USER_SYSTEM,
		To:        sender.From,
		Content:   sender.Id, // 后端分配好的消息ID
		Type:      constant.SYSTEM_INFO,
		MsgType:   constant.MSG_ACK_MSG,
		Timestamp: time.Now().Unix(),
	}
	protobuf, err := proto.Marshal(&reply)
	if err != nil {
		logx.Error("[replyAckMessage] Marshal message failed, error: ", err)
		return
	}
	l.packageMessageAndSend(protobuf, sender.From, oldId, constant.MSG_ACK_MSG)
}

func (l *MsgForwarder) systemOperation(message *front.Message, protobuf []byte) {
	l.packageMessageAndSend(protobuf, message.To, "", message.MsgType)
}
