package svc

import (
	"github.com/n8sPxD/cowIM/microservices/wsget/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
