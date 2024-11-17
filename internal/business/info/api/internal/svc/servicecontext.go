package svc

import (
	"github.com/n8sPxD/cowIM/internal/business/info/api/internal/config"
	"github.com/n8sPxD/cowIM/internal/common/db/myMongo"
)

type ServiceContext struct {
	Config config.Config
	Mongo  *myMongo.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Mongo:  myMongo.MustNewMongo(c.Mongo.Host),
	}
}
