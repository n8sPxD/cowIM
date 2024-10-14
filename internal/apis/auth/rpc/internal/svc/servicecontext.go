package svc

import (
	"github.com/n8sPxD/cowIM/internal/apis/auth/rpc/internal/config"
	"github.com/n8sPxD/cowIM/internal/common/db/myRedis"
)

type ServiceContext struct {
	Config config.Config
	Redis  *myRedis.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  myRedis.MustNewRedis(c.RedisConf),
	}
}
