package svc

import (
	"time"

	"github.com/n8sPxD/cowIM/common/db"
	"github.com/n8sPxD/cowIM/imrpc/internal/config"
	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config         config.Config
	MySQL          *db.DB
	KqPusherClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqldb := db.MustNewMySQL(c.MySQL.DataSource)
	return &ServiceContext{
		Config:         c,
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic, kq.WithFlushInterval(time.Millisecond*500)),
		MySQL:          mysqldb,
	}
}
