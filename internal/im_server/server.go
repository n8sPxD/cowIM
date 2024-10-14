package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/n8sPxD/cowIM/internal/im_server/internal/config"
	"github.com/n8sPxD/cowIM/internal/im_server/internal/mqs"
	"github.com/n8sPxD/cowIM/internal/im_server/internal/server"
	"github.com/n8sPxD/cowIM/internal/im_server/svc"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

var configFile = flag.String("f", "etc/server.yaml", "the config file")

func main() {
	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	idOptions := idgen.NewIdGeneratorOptions(c.WorkID)
	idgen.SetIdGenerator(idOptions)

	ctx := context.Background()
	svcCtx := svc.NewServiceContext(c)

	var mq *mqs.MsgSender
	threading.GoSafe(func() {
		mq = mqs.NewMsgSender(ctx, svcCtx)
		mq.Start()
	})

	// 处理退出信号，平滑关闭
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		mq.Close()
		os.Exit(0)
	}() // 处理退出信号，平滑关闭

	s := server.MustNewServer(ctx, svcCtx)
	s.Start()
}
