package main

import (
	"flag"
	"fmt"

	"github.com/n8sPxD/cowIM/internal/business/auth/api/internal/config"
	"github.com/n8sPxD/cowIM/internal/business/auth/api/internal/handler"
	"github.com/n8sPxD/cowIM/internal/business/auth/api/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/auth.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting socket at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
