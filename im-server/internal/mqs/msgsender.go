package mqs

import (
	"context"

	__front "github.com/n8sPxD/cowIM/common/message/.front"
	"github.com/n8sPxD/cowIM/im-server/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgSender struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMsgSender(ctx context.Context, svcCtx *svc.ServiceContext) *MsgSender {
	return &MsgSender{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Consume 接受从 后端消息处理服务 发来的消息，并转发给对应的用户
func (l *MsgSender) Consume(ctx context.Context, key, val string) error {
	// manager := l.svcCtx.ConnectionManager
	logx.Infof("PaymentSuccess key :%s , val :%s", key, val)
	l.svcCtx.ConnectionManager.SendMessage(7, &__front.Message{
		Content: "哈哈",
	})
	return nil
}
