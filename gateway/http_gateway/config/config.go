package config

import "github.com/zeromicro/go-zero/core/logx"

type Proxy struct {
	Route  string
	Target string
}

type Config struct {
	Proxies []Proxy
	Log     logx.LogConf
}
