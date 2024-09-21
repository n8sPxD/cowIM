package svc

import (
	"github.com/n8sPxD/cowIM/common/db/myMysql"
	"github.com/n8sPxD/cowIM/microservices/auth/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
	MySQL  *myMysql.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqldb := myMysql.MustNewMySQL(c.MySQL.DataSource)
	return &ServiceContext{
		Config: c,
		MySQL:  mysqldb,
	}
}
