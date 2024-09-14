package main

import (
	"flag"
	"fmt"

	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/config"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/server"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	svcCtx := svc.NewServiceContext(c)

	tcps := server.MustNewTcpServer(svcCtx, "0.0.0.0:8000")

	fmt.Printf("Starting tcp server at %s...\n", "0.0.0.0:8000")

	tcps.HandleRequest()
}
