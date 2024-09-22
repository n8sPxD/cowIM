package mqs

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/message/front"
	"github.com/n8sPxD/cowIM/im-server/internal"
	"github.com/n8sPxD/cowIM/im-server/svc"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type MsgSender struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	MsgSender *kafka.Reader
}

func NewMsgSender(ctx context.Context, svcCtx *svc.ServiceContext) *MsgSender {
	return &MsgSender{
		ctx:    ctx,
		svcCtx: svcCtx,
		MsgSender: kafka.NewReader(kafka.ReaderConfig{
			Brokers: svcCtx.Config.MsgSender.Brokers,
			Topic:   svcCtx.Config.MsgSender.Topic,
		}),
	}
}

func (l *MsgSender) Start() {
	// 设置kafka起始偏移量，在初始化NewReader的时候设置没用不知道为什么，只有这里有用
	err := l.MsgSender.SetOffset(kafka.LastOffset)
	if err != nil {
		logx.Error("[MsgForwarder.Start] Set kafka offset failed, error: ", err)
	}

	for {
		msg, err := l.MsgSender.ReadMessage(l.ctx)
		if err != nil {
			logx.Error("[MsgSender.Start] Reading msgForward error: ", err)
			continue
		}
		go l.Consume(msg.Value)
	}
}

func (l *MsgSender) Close() {
	_ = l.MsgSender.Close()
}

// Consume 接受从 后端消息处理服务 发来的消息，并转发给对应的用户
func (l *MsgSender) Consume(protobuf []byte) {
	var msg front.Message
	err := proto.Unmarshal(protobuf, &msg)
	if err != nil {
		logx.Error("[MsgSender.Consume] Protobuf unmarshal failed, error: ", err)
		return
	}
	switch msg.Type {
	case constant.SINGLE_CHAT:
		l.SingleChat(&msg, protobuf)
	case constant.GROUP_CHAT:
		l.GroupChat(&msg)
	case constant.BIG_GROUP_CHAT:
		l.BigGroupChat()
	default:
		logx.Error("[MsgSender.Consume] Wrong msgForward type, Type is: ", msg.Type)
	}
}

func (l *MsgSender) SingleChat(msg *front.Message, protobuf []byte) {
	// 能传到这里来，代表Message服务中已经从Redis中获取到当前Recv用户在线
	// 在线，以服务器主动推模式发送消息
	err := l.svcCtx.ConnectionManager.SendMessage(msg.To, protobuf)
	if err != nil {
		// Message服务中检测到用户在线，但是可能在消息中转的过程中又离线
		if errors.Is(err, internal.ClientGoingAway) {
			// TODO: 更改Redis信息
			return
		}
		logx.Error("[MsgSender.SingleChat] Send msgForward failed, error: ", err)
		return
	}
}

func (l *MsgSender) GroupChat(msg *front.Message) {
	// TODO: 完善逻辑
}

func (l *MsgSender) BigGroupChat() {
	// TODO: 完善逻辑
}
