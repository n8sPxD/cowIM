package svc

import (
	"fmt"

	"github.com/n8sPxD/cowIM/internal/common/db/myMongo"
	"github.com/n8sPxD/cowIM/internal/msg_to_db/internal/config"
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
