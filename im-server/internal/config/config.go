package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int
	}
	WorkID        uint16
	AuthRpc       zrpc.RpcClientConf
	MsgPusherConf struct {
		Brokers []string
		Topic   string
	}
	SendConsumerConf kq.KqConf
}
