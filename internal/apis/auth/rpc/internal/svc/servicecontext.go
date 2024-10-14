package svc

import (
	"github.com/n8sPxD/cowIM/internal/apis/auth/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.MustNewRedis(c.RedisConf)
	return &ServiceContext{
		Config: c,
		Redis:  rdb,
	}
}
