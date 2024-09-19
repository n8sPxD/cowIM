package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/n8sPxD/cowIM/im-server/internal/config"
	"github.com/n8sPxD/cowIM/im-server/internal/mqs"
	"github.com/n8sPxD/cowIM/im-server/internal/server"
	"github.com/n8sPxD/cowIM/im-server/svc"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/threading"
)

var configFile = flag.String("f", "etc/server.yaml", "the config file")

func main() {
	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	idOptions := idgen.NewIdGeneratorOptions(c.WorkID)
	idgen.SetIdGenerator(idOptions)

	// 处理退出信号，平滑关闭
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		os.Exit(0)
	}() // 处理退出信号，平滑关闭

	ctx := context.Background()
	svcCtx := svc.NewServiceContext(c)

	threading.GoSafe(func() {
		s := server.MustNewServer(c, ctx, svcCtx)
		svcCtx.ConnectionManager = s.Manager
		s.Start()
	})

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	for _, mq := range mqs.Consumers(c, ctx, svcCtx) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()
}
