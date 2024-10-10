package svc

import (
	"time"

	"github.com/n8sPxD/cowIM/common/db/myRedis"
	"github.com/n8sPxD/cowIM/common/servicehub"
	"github.com/n8sPxD/cowIM/im-server/internal/config"
	"github.com/n8sPxD/cowIM/im-server/internal/server/manager"
	"github.com/n8sPxD/cowIM/microservices/auth/rpc/auth"
	"github.com/n8sPxD/cowIM/microservices/auth/rpc/types/authRpc"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	AuthRpc           authRpc.AuthClient
	ConnectionManager *manager.ConnectionManager
	MsgForwarder      *kafka.Writer
	Redis             *myRedis.DB
	RegisterHub       *servicehub.RegisterHub
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		AuthRpc: auth.NewAuth(zrpc.MustNewClient(c.AuthRpc)),
		MsgForwarder: &kafka.Writer{
			Addr:         kafka.TCP(c.MsgForwarder.Brokers...),
			Topic:        c.MsgForwarder.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond, // 低超时时间
			RequiredAcks: kafka.RequireOne,      // 仅等待 Leader 确认
			Compression:  compress.Zstd,         // Zstd压缩
			Async:        true,                  // 启用异步写入
			MaxAttempts:  1,                     // 限制重试次数
		},
		Redis:       myRedis.MustNewRedis(c.RedisConf),
		RegisterHub: servicehub.NewRegisterHub(c.Etcd.Hosts, 3),
	}
}
