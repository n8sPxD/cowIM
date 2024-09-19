package svc

import (
	"github.com/n8sPxD/cowIM/common/db/mysql"
	"github.com/n8sPxD/cowIM/microservices/auth/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
	MySQL  *mysql.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqldb := mysql.MustNewMySQL(c.MySQL.DataSource)
	return &ServiceContext{
		Config: c,
		MySQL:  mysqldb,
	}
}
