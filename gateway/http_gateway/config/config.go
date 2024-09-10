package config

import "github.com/zeromicro/go-zero/core/logx"

type Proxy struct {
	Route  string
	Target string
}

type Config struct {
	Log     logx.LogConf
	Proxies []Proxy
	Auth    struct {
		AccessSecret string
		AccessExpire int
	}
	WhiteList []string
}
