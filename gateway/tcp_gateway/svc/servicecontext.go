package svc

import (
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
