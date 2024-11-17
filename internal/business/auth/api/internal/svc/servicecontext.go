package svc

import (
	"github.com/n8sPxD/cowIM/internal/business/auth/api/internal/config"
	"github.com/n8sPxD/cowIM/internal/common/db/myMysql"
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
