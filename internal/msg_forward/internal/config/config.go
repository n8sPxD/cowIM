package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	Log       logx.LogConf
	RedisConf redis.RedisConf
	MongoConf struct {
		Host string
	}
	MsgForwarder struct {
		Brokers []string
		Topic   string
	}
	MsgSender struct {
		Brokers []string
		Topic   string
	}
	MsgDBSaver struct {
		Brokers []string
		Topic   string
	}
	MySQL struct {
		DataSource string
	}
	WorkID uint16
}
