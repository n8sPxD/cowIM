package main

import (
	"flag"
	"fmt"

	"github.com/n8sPxD/cowIM/internal/business/group/rpc/internal/config"
	"github.com/n8sPxD/cowIM/internal/business/group/rpc/internal/server"
	"github.com/n8sPxD/cowIM/internal/business/group/rpc/internal/svc"
	"github.com/n8sPxD/cowIM/internal/business/group/rpc/types/group"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/group.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		group.RegisterGroupServer(grpcServer, server.NewGroupServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
