package svc

import (
	"github.com/n8sPxD/cowIM/common/db/myMysql"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
    MySQL *myMysql.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
        MySQL: myMysql.MustNewMySQL(c.MySQL.DataSource),
	}
}
