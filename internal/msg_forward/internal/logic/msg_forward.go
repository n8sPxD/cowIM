package logic

import (
	"context"
	"github.com/n8sPxD/cowIM/internal/msg_forward/gossip"
	"github.com/segmentio/kafka-go/compress"
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

// 群聊处理
func (l *MsgForwarder) groupChat(msg *front.Message, protobuf []byte) {
	// 先获取群里所有成员
	// TODO: 群组成员加缓存
	members, err := l.svcCtx.MySQL.GetGroupMemberIDs(l.ctx, uint(*msg.Group))
	if err != nil {
		logx.Error("[groupChat] Get group members from mysql failed, error: ", err)
		return
	}
	// TODO: 可以优化，先处理所有消息，然后把对应服务器的消息以切片形式发送，避免重复调用WriteMessages
	for _, member := range members {
		// 不用给发消息的人发
		if member == uint(msg.From) {
			continue
		}
		// TODO: Redis中用Pipeline一次性查询所有member对应的服务器ID, 直接转发，避免重复的单次查询
		l.packageMessageAndSend(protobuf, msg.To, msg.Id, msg.MsgType)
	}
}

/*
大型群聊 （大于500人） 客户端定期轮询请求
客户端登陆，同步最新Timeline，服务端不主动推送大型群聊消息
客户端定期轮询（1秒），发送请求，服务端返回最新的一条消息
当客户端点进群聊，发送请求，服务端返回服务端最新消息和客户端最新消息之间的所有消息

客户端刚登陆的时候，假如加入了大规模群组1001，则定期给服务器发送消息：

	{
		from = A
		to = SYSTEM
		group = 1001
		content = "所有群组的id,用下划线_拼接起来"
		type = BIG_GROUP
		msg_type = BIG_GROUP_REQ
		timestamp = 当前本地所有BIG_GROUP消息内的最新时间戳
	}

此时服务器收到消息，将当前本地（缓存或数据库）的最新一条记录发过去：

	{
		from = SYSTEM
		to =  A
		group = 1001
		content = 群聊1001的最新一条消息
		type = BIG_GROUP
		msg_type = BIG_GROUP_REQ
		timestamp = timestamp
	}

客户端收到后，在前端将最新消息显示出来，但是消息不入库，在客户端在线的时间持续间接性发送。
当客户端具体点击到一个群之后，前端进行完整的定时轮询，每次向服务器发送消息如下：

	{
		from = A
		to = SYSTEM
		group = 1001
		content = ""
		type = BIG_GROUP
		msg_type = BIG_GROUP_ALL_REQ
		timestamp = 当前本地所有BIG_GROUP消息内的最新时间戳
	}

服务器接收到了，然后返回消息：

	{
		id = 消息本身的id
		from = SYSTEM
		to = A
		group = 1001
		content = 消息内容
		type = BIG_GROUP
		msg_type = MSG_COMMON_MSG
	}

前端接收到，返回Ack，存数据库，如果有其他群例如1002,那么一样的流程。
该流程可以优化，例如第一次握手时，客户端以切片形式发送请求
*/
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
