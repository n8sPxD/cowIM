package svc

import (
	"github.com/n8sPxD/cowIM/common/db"
	"github.com/n8sPxD/cowIM/imapi/auth/internal/config"
)

type ServiceContext struct {
	Config config.Config
	MySQL  *db.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqldb := db.MustNewMySQL(c.MySQL.DataSource)
	return &ServiceContext{
		Config: c,
		MySQL:  mysqldb,
	}
}
