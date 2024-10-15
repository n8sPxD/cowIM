package mqs

import (
	"context"
	"time"

	"github.com/n8sPxD/cowIM/internal/common/constant"
	"github.com/n8sPxD/cowIM/internal/common/message/inside"
	"github.com/n8sPxD/cowIM/internal/im_server/internal/server"
	"github.com/n8sPxD/cowIM/internal/im_server/svc"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type MsgSender struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	manager   server.IConnectionManager
	MsgSender *kafka.Reader
}

func NewMsgSender(ctx context.Context, svcCtx *svc.ServiceContext) *MsgSender {
	return &MsgSender{
		ctx:    ctx,
		svcCtx: svcCtx,
		MsgSender: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        svcCtx.Config.MsgSender.Brokers,
			Topic:          svcCtx.Config.MsgSender.Topic,
			GroupID:        "msg-fwd",
			StartOffset:    kafka.LastOffset,
			MinBytes:       1,                      // 最小拉取字节数
			MaxBytes:       10e3,                   // 最大拉取字节数（10KB）
			MaxWait:        100 * time.Millisecond, // 最大等待时间
			CommitInterval: 500 * time.Millisecond, // 提交间隔
		}),
	}
}

func (l *MsgSender) WithManager(manager server.IConnectionManager) *MsgSender {
	l.manager = manager
	return l
}

func (l *MsgSender) Start() {
	for {
		msg, err := l.MsgSender.ReadMessage(l.ctx)
		if err != nil {
			logx.Error("[MsgSender.Start] Reading message error: ", err)
			continue
		}
		go l.Consume(msg.Value)
	}
}

func (l *MsgSender) Close() {
	l.MsgSender.Close()
}

// Consume 接受从 后端消息处理服务 发来的消息，并转发给对应的用户
func (l *MsgSender) Consume(protobuf []byte) {
	var msg inside.Message
	if err := proto.Unmarshal(protobuf, &msg); err != nil {
		logx.Error("[MsgSender.Consume] Protobuf unmarshal failed, error: ", err)
		return
	}
	// 能传到这里来，代表Message服务中已经从Redis中获取到当前Recv用户在线
	// 在线，以服务器主动推模式发送消息
	switch msg.Type {
	case constant.MSG_ACK_MSG:
		// 有用户发送消息到服务器，需要获得来自服务器的Ack
		go l.replyMessageWithAck(&msg)
	default:
		// 一般正常消息
		go l.sendMessage(&msg, 2*time.Second, 3)
	}
}

func (l *MsgSender) replyMessageWithAck(message *inside.Message) {
	logx.Debug("[replyMessageWithAck] Replying ack message to User ", message.To)
	if err := l.manager.SendMessage(message.To, message.Protobuf); err != nil {
		logx.Error("[replyMessageWithAck] Reply Ack message failed, error: ", err)
	}
}

func (l *MsgSender) sendMessage(message *inside.Message, retryInterval time.Duration, maxRetires int) {
	for range maxRetires {
		var (
			ackChan    = make(chan bool)
			manager    = l.manager
			ackHandler = manager.GetAckHandler()
			ack        = server.Ack{To: message.To, MessageID: message.MsgId}
		)
		// 创建等待用Ack Channel
		ackHandler.AssignAckChan(ack, ackChan)
		// 等待Ack或超时
		go ackHandler.WaitForAck(ack, retryInterval)
		// 发送消息
		if err := manager.SendMessage(message.To, message.Protobuf); err != nil {
			logx.Errorf("[sendMessage] Send message to User %d failed, error: ", err)
		}
		// 等待Ack
		if <-ackChan {
			logx.Debugf("[sendMessage] Receive Ack from User %d with message \"%s\"", message.To, message.MsgId)
			return
		}
	}
	logx.Errorf("[sendMessage] Receive Ack from User %d with message \"%s\" failed", message.To, message.MsgId)
}
