package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Log  logx.LogConf
	Auth struct {
		AccessSecret string
		AccessExpire int
	}
	IMRpc          zrpc.RpcClientConf
	KqConsumerConf kq.KqConf
}
