package main

import (
	"context"
	"flag"

	"github.com/n8sPxD/cowIM/microservices/message/internal/config"
	"github.com/n8sPxD/cowIM/microservices/message/internal/mqs"
	"github.com/n8sPxD/cowIM/microservices/message/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/message.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()

	mq := mqs.NewMsgForwarder(ctx, svcCtx)
	mq.Start()
}
