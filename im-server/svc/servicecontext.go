package svc

import (
	"github.com/n8sPxD/cowIM/im-server/internal"
	"github.com/n8sPxD/cowIM/im-server/internal/config"
	"github.com/n8sPxD/cowIM/microservices/auth/rpc/auth"
	"github.com/n8sPxD/cowIM/microservices/auth/rpc/types/authRpc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	AuthRpc           authRpc.AuthClient
	MsgPusher         *kq.Pusher
	ConnectionManager *internal.ConnectionManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		AuthRpc: auth.NewAuth(zrpc.MustNewClient(c.AuthRpc)),
		MsgPusher: kq.NewPusher(
			c.MsgPusherConf.Brokers,
			c.MsgPusherConf.Topic, // TODO: 限制chunkSize等
		),
	}
}
