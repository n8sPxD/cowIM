package main

import (
	"context"
	"flag"

	"github.com/n8sPxD/cowIM/microservices/msgToDB/internal/config"
	"github.com/n8sPxD/cowIM/microservices/msgToDB/internal/mqs"
	"github.com/n8sPxD/cowIM/microservices/msgToDB/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()

	mq := mqs.NewMsgToDB(ctx, svcCtx)
	mq.Start()
}
