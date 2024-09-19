package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
)

type Config struct {
	Log             logx.LogConf
	MsgConsumerConf kq.KqConf
	SendPusherConf  struct {
		Brokers []string
		Topic   string
	}
}
