package main

import (
	"context"
	"flag"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/config"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/logic"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"os/signal"
	"syscall"
)

var (
	configFile = flag.String("f", "etc/message.yaml", "the config file")

	c  config.Config
	mq *logic.MsgForwarder
)

func main() {
	flag.Parse()

	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()

	mq = logic.NewMsgForwarder(ctx, svcCtx)
	go mq.Start()

	svcCtx.Regist.Register(ctx, svcCtx.Config.Name, svcCtx.Config.Port, svcCtx.Config.WorkID)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	svcCtx.Redis.RemoveAllDupMessages(ctx) // 移除所有未处理的重复消息uuid
	mq.Close()
	os.Exit(0)
}
