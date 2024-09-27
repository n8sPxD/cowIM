package svc

import (
	"github.com/n8sPxD/cowIM/common/db/myMongo"
	"github.com/n8sPxD/cowIM/common/db/myMysql"
	"github.com/n8sPxD/cowIM/microservices/info/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Mongo  *myMongo.DB
	MySQL  *myMysql.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Mongo:  myMongo.MustNewMongo(c.Mongo.Host),
		MySQL:  myMysql.MustNewMySQL(c.MySQL.DataSource),
	}
}
