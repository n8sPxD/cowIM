package svc

import (
	"github.com/n8sPxD/cowIM/common/db/myMongo"
	"github.com/n8sPxD/cowIM/microservices/info/internal/config"
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
