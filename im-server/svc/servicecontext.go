package svc

import (
	"github.com/n8sPxD/cowIM/im-server/internal"
	"github.com/n8sPxD/cowIM/im-server/internal/config"
	"github.com/n8sPxD/cowIM/microservices/auth/rpc/auth"
	"github.com/n8sPxD/cowIM/microservices/auth/rpc/types/authRpc"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	AuthRpc           authRpc.AuthClient
	ConnectionManager *internal.ConnectionManager
	MsgForwarder      *kafka.Writer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		AuthRpc: auth.NewAuth(zrpc.MustNewClient(c.AuthRpc)),
		MsgForwarder: &kafka.Writer{
			Addr:  kafka.TCP(c.MsgForwarder.Brokers...),
			Topic: c.MsgForwarder.Topic,
		},
	}
}
