package mqs

import (
	"context"

	"github.com/n8sPxD/cowIM/im-server/internal/config"
	"github.com/n8sPxD/cowIM/im-server/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func Consumers(
	c config.Config,
	ctx context.Context,
	svcContext *svc.ServiceContext,
) []service.Service {
	return []service.Service{
		kq.MustNewQueue(c.SendConsumerConf, NewMsgSender(ctx, svcContext)),
	}
}
