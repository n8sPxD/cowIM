package svc

import (
	"fmt"

	"github.com/n8sPxD/cowIM/common/db/myMongo"
	"github.com/n8sPxD/cowIM/microservices/msgToDB/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Mongo  *myMongo.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Mongo:  myMongo.MustNewMongo(fmt.Sprintf("mongodb://%s", c.MongoConf.Host)),
	}
}
