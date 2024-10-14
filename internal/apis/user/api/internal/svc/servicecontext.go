package svc

import (
	"github.com/n8sPxD/cowIM/internal/apis/user/api/internal/config"
	"github.com/n8sPxD/cowIM/internal/common/db/myMysql"
)

type ServiceContext struct {
	Config config.Config
	MySQL  *myMysql.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		MySQL:  myMysql.MustNewMySQL(c.MySQL.DataSource),
	}
}
