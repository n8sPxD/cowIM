package svc

import (
	"github.com/n8sPxD/cowIM/common/servicehub"
	"github.com/n8sPxD/cowIM/microservices/wsget/internal/config"
)

type ServiceContext struct {
	Config       config.Config
	DiscoveryHub *servicehub.DiscoveryHub
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		DiscoveryHub: servicehub.NewDiscoveryHub(c.Etcd.Host, 3),
	}
}
