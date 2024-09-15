package svc

import (
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/config"
	"github.com/n8sPxD/cowIM/imrpc/imrpcclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	IMRpc  imrpcclient.ImRPC
}

func NewServiceContext(c config.Config) *ServiceContext {
	rpc := imrpcclient.NewImRPC(zrpc.MustNewClient(c.IMRpc))
	return &ServiceContext{
		Config: c,
		IMRpc:  rpc,
	}
}
