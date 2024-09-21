package mqs

import (
	"context"
	"fmt"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/message/front"
	"github.com/n8sPxD/cowIM/microservices/message/internal/svc"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
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
			Brokers: svcCtx.Config.MsgForwarder.Brokers,
			Topic:   svcCtx.Config.MsgForwarder.Topic,
		}),
	}
}

func (l *MsgForwarder) Start() {
	for {
		msg, err := l.MsgForwarder.ReadMessage(l.ctx) // 这里的msg是kafka.Message
		if err != nil {
			logx.Error("[MsgForwarder.Start] Reading message error: ", err)
			continue
		}
		logx.Info("value: ", msg)
		go l.Consume(msg.Value)
	}
}

// Consume 接收从 Websocket Server的消息，处理后再进行转发
func (l *MsgForwarder) Consume(protobuf []byte) {
	// 传过来的消息是序列化过的，先反序列化
	var msg front.Message
	err := proto.Unmarshal(protobuf, &msg)
	if err != nil {
		logx.Error("[MsgForwarder.Consume] Unmarshal message failed, error: ", err)
		return
	}

	// TODO: 把消息存MessageRecord表中

	switch msg.Type {
	case constant.SINGLE_CHAT:
		go l.SingleChat(&msg, protobuf)
	case constant.GROUP_CHAT:
		go l.GroupChat(&msg)
	case constant.BIG_GROUP_CHAT:
		go l.BigGroupChat(&msg)
	default:
		logx.Error("[MsgForwarder.Consume] Wrong message type, Type is: ", msg.Type)
	}
}

func (l *MsgForwarder) SingleChat(msg *front.Message, protobuf []byte) {
	// 查询Redis中路由信息
	status, err := l.svcCtx.Redis.GetUserRouterStatus(l.ctx, msg.To)
	if err != nil {
		logx.Error("[MsgForwarder.SingleChat] Get router status from redis failed, error: ", err)
		return
	}
	if status == nil {
		// 没找到当前用户的路由信息，说明没上线
		// TODO: 塞timeline里
		return
	}
	// 转发消息到指定的websocket-server
	// 先确定Topic
	workID := status.WorkID
	l.svcCtx.MsgSender.Topic = fmt.Sprintf("websocket-server-%d", workID)
	// 封装消息
	mqMsg := kafka.Message{
		Value: protobuf,
	}
	err = l.svcCtx.MsgSender.WriteMessages(l.ctx, mqMsg)
	if err != nil {
		logx.Error("[MsgForwarder.SingleChat] Push message to MQ failed, error: ", err)
		return
	}
}

func (l *MsgForwarder) GroupChat(msg *front.Message) {
	// TODO: 完善逻辑
}

func (l *MsgForwarder) BigGroupChat(msg *front.Message) {
	// TODO: 完善逻辑
}
