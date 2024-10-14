package mqs

import (
	"context"
	"errors"
	"time"

	"github.com/n8sPxD/cowIM/internal/common/message/inside"
	"github.com/n8sPxD/cowIM/internal/im_server/internal/server/manager"
	"github.com/n8sPxD/cowIM/internal/im_server/svc"
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

func (l *MsgSender) Start() {
	for {
		msg, err := l.MsgSender.ReadMessage(l.ctx)
		if err != nil {
			logx.Error("[MsgSender.Start] Reading message error: ", err)
			continue
		}
		logx.Debugf(
			"[MsgForwarder.Start] Message at partition %d offset %d: %s\n",
			msg.Partition,
			msg.Offset,
			string(msg.Value),
		)
		go l.Consume(msg.Value)
	}
}

func (l *MsgSender) Close() {
	_ = l.MsgSender.Close()
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
	if err := l.svcCtx.ConnectionManager.SendMessage(msg.To, msg.Protobuf); err != nil {
		// Message服务中检测到用户在线，但是可能在消息中转的过程中又离线
		if errors.Is(err, manager.ClientGoingAway) {
			// TODO: 更改Redis信息
			return
		}
		logx.Error("[MsgSender.SingleChat] Send message failed, error: ", err)
		return
	}
}
