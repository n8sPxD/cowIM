package mqs

import (
	"context"

	"github.com/n8sPxD/cowIM/microservices/message/internal/config"
	"github.com/n8sPxD/cowIM/microservices/message/internal/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func Consumers(
	c config.Config,
	ctx context.Context,
	svcContext *svc.ServiceContext,
) []service.Service {
	return []service.Service{
		// Listening for changes in consumption flow status
		kq.MustNewQueue(c.MsgConsumerConf, NewMsgConsumer(ctx, svcContext)),
	}
}
