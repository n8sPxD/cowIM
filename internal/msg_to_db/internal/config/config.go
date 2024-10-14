package config

import "github.com/zeromicro/go-zero/core/logx"

type Config struct {
	Log       logx.LogConf
	MongoConf struct {
		Host string
	}
	MsgToDB struct {
		Brokers []string
		Topic   string
	}
}
