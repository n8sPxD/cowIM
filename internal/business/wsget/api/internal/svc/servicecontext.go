package svc

import (
	"github.com/n8sPxD/cowIM/internal/business/wsget/api/internal/config"
	"github.com/n8sPxD/cowIM/pkg/servicehub"
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
