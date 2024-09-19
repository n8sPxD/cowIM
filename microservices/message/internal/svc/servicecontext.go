package svc

import (
	"github.com/n8sPxD/cowIM/microservices/message/internal/config"
	"github.com/zeromicro/go-queue/kq"
)

type ServiceContext struct {
	Config     config.Config
	SendPusher *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		SendPusher: kq.NewPusher(c.SendPusherConf.Brokers, c.SendPusherConf.Topic),
	}
}
