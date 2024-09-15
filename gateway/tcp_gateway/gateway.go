package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/config"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/mqs"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/server"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/threading"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()

	tcps := server.MustNewTcpServer(svcCtx, "0.0.0.0:8000")

	fmt.Printf("Starting tcp server at %s...\n", "0.0.0.0:8000")

	threading.GoSafe(func() { tcps.HandleRequest() })

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range mqs.Consumers(c, ctx, svcCtx, tcps.Server) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()
}
