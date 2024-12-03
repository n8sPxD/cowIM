package svc

import (
	"github.com/n8sPxD/cowIM/internal/business/group/rpc/internal/config"
	"github.com/n8sPxD/cowIM/internal/common/dao/myMysql"
)

type ServiceContext struct {
	Config config.Config
	DB     *myMysql.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     myMysql.MustNewMySQL(c.MySQL.DataSource),
	}
}
