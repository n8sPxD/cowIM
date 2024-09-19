package mqs

import (
	"context"
	"fmt"

	"github.com/n8sPxD/cowIM/microservices/message/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgConsumer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMsgConsumer(ctx context.Context, svcCtx *svc.ServiceContext) *MsgConsumer {
	return &MsgConsumer{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Consume 接收从 Websocket Server的消息，处理后再进行转发
func (l *MsgConsumer) Consume(ctx context.Context, key, val string) error {
	logx.Infof("PaymentSuccess key :%s , val :%s", key, val)
	if err := l.svcCtx.SendPusher.Push(ctx, fmt.Sprintf("%s:%s", key, val)); err != nil {
		logx.Error("err! ", err)
	}
	return nil
}
