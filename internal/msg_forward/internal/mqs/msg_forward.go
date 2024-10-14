package mqs

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/n8sPxD/cowIM/internal/common/constant"
	"github.com/n8sPxD/cowIM/internal/common/message/front"
	"github.com/n8sPxD/cowIM/internal/common/message/inside"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/svc"
	"github.com/segmentio/kafka-go"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"google.golang.org/protobuf/proto"
)

type MsgForwarder struct {
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	MsgForwarder *kafka.Reader
}

func NewMsgForwarder(ctx context.Context, svcCtx *svc.ServiceContext) *MsgForwarder {
	return &MsgForwarder{
		ctx:    ctx,
		svcCtx: svcCtx,
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
	}
}

func (l *MsgForwarder) Close() {
	l.MsgForwarder.Close()
}

func (l *MsgForwarder) Start() {
	// 初始化id生成器
	options := idgen.NewIdGeneratorOptions(l.svcCtx.Config.WorkID)
	idgen.SetIdGenerator(options)

	for {
		msg, err := l.MsgForwarder.ReadMessage(l.ctx) // 这里的msg是kafka.Message
		if err != nil {
			logx.Error("[MsgForwarder.Start] Reading message error: ", err)
			break
		}
		go l.Consume(msg.Value, time.Now())
	}
}

// Consume 接收从 Websocket Server的消息，处理后再进行转发
func (l *MsgForwarder) Consume(protobuf []byte, now time.Time) {
	// 传过来的消息是序列化过的，先反序列化
	var msg front.Message
	err := proto.Unmarshal(protobuf, &msg)
	if err != nil {
		logx.Error("[Consume] Unmarshal message failed, error: ", err)
		return
	}

	// 检查消息重复性
	if ok, err := l.svcCtx.Redis.CheckDuplicateMessage(l.ctx, msg.Id); err != nil {
		logx.Error("[Consume] Check duplicate message from redis failed, error: ", err)
		return
	} else if ok { // 消息是重复的
		logx.Infof("[Consume] Message from %d with ID \"%s\" sent repeated", msg.From, msg.Id)
		return
	}

	// 保存以前的uuid
	oldId := msg.Id
	// 分配消息ID
	id := idgen.NextId()
	msg.Id = strconv.FormatInt(id, 10)

	protobuf, err = proto.Marshal(&msg)

	// 异步存库
	go l.sendRecordMsgToDB(&msg, now) // 漫游库
	go l.sendTimelineToDB(&msg, now)  // timeline

	// Ack消息
	go l.replyAckMessage(&msg, oldId)

	// 进行基于消息类型的消息处理
	switch msg.Type {
	case constant.SINGLE_CHAT:
		go l.singleChat(&msg, protobuf)
	case constant.GROUP_CHAT:
		go l.groupChat(&msg, protobuf)
	case constant.BIG_GROUP_CHAT:
		go l.bigGroupChat(&msg)
	default:
		logx.Error("[Consume] Wrong message type, Type is: ", msg.Type)
		return
	}
}

func (l *MsgForwarder) packageMessage(protobuf []byte, id uint32) (kafka.Message, error) {
	// 先查用户在不在线
	status, err := l.svcCtx.Redis.GetUserRouterStatus(l.ctx, id)
	if errors.Is(err, redis.Nil) {
		// 不在线，直接跳过
		return kafka.Message{}, redis.Nil
	} else if err != nil {
		// redis出问题了
		return kafka.Message{}, err
	}

	// 用户在线，发消息
	// 确定Topic
	workID := status.WorkID
	l.svcCtx.MsgSender.Topic = fmt.Sprintf("websocket-server-%d", workID)

	// 封装inside.Message
	msg := inside.Message{
		To:       id,
		Protobuf: protobuf,
	}
	msgByte, err := proto.Marshal(&msg)
	if err != nil {
		return kafka.Message{}, err
	}

	return kafka.Message{Value: msgByte}, nil
}

// 单聊处理
func (l *MsgForwarder) singleChat(msg *front.Message, protobuf []byte) {
	km, err := l.packageMessage(protobuf, msg.To)
	if err != nil {
		logx.Error("[singleChat] Package message failed, error: ", err)
		return
	}
	err = l.svcCtx.MsgSender.WriteMessages(l.ctx, km)
	if err != nil {
		logx.Error(
			"[singleChat] Push message to Websocket-server MQ failed, error: ",
			err,
		)
		return
	}
}

// 群聊处理
func (l *MsgForwarder) groupChat(msg *front.Message, protobuf []byte) {
	// 先获取群里所有成员
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
		kq, err := l.packageMessage(protobuf, uint32(member))
		if err != nil {
			logx.Error("[groupChat] Package message failed, error: ", err)
			return
		}
		// 发消息
		if err := l.svcCtx.MsgSender.WriteMessages(l.ctx, kq); err != nil {
			logx.Error("[groupChat] Push message to Websocket-server MQ failed, error: ", err)
			continue
		}
	}
}

func (l *MsgForwarder) bigGroupChat(msg *front.Message) {
	// TODO: 完善逻辑
}

func (l *MsgForwarder) replyAckMessage(sender *front.Message, oldId string) {
	// 封装消息体
	reply := front.Message{
		Id:      oldId,
		From:    constant.USER_SYSTEM,
		To:      sender.From,
		Content: sender.Id, // 后端分配好的消息ID
		Type:    constant.SYSTEM_INFO,
		MsgType: constant.MSG_ACK_MSG,
	}
	protobuf, err := proto.Marshal(&reply)
	if err != nil {
		logx.Error("[replyAckMessage] Marshal message failed, error: ", err)
		return
	}
	kq, err := l.packageMessage(protobuf, sender.From)
	if err != nil {
		logx.Error("[replyAckMessage] Package message failed, error: ", err)
		return
	}
	if err := l.svcCtx.MsgSender.WriteMessages(l.ctx, kq); err != nil {
		logx.Error("[replyAckMessage] Push message to Websocket-server MQ failed, error: ", err)
		return
	}
}
