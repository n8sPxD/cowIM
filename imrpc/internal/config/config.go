package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MySQL struct {
		DataSource string
	}
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
}
