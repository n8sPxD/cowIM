package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/config"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/mqs"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	configFile = flag.String("f", "etc/message.yaml", "the config file")

	c  config.Config
	mq *mqs.MsgForwarder
)

func main() {
	flag.Parse()

	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()

	// 处理退出信号，平滑关闭
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		svcCtx.Redis.RemoveAllDupMessages() // 移除所有未处理的重复消息uuid
		mq.Close()
		os.Exit(0)
	}() // 处理退出信号，平滑关闭

	mq = mqs.NewMsgForwarder(ctx, svcCtx)
	mq.Start()
}
