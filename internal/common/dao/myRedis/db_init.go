package myRedis

import (
	"github.com/redis/go-redis/v9"
	redis2 "github.com/zeromicro/go-zero/core/stores/redis"
	"strings"
)

type Native struct {
	redis.UniversalClient
}

type GoZero struct {
	*redis2.Redis
}

func MustNewNativeRedis(c redis2.RedisConf) *Native {
	return &Native{redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    strings.Split(c.Host, ","),
		Password: c.Pass,
	})}
}

func MustNewGoRedis(c redis2.RedisConf) *GoZero {
	return &GoZero{redis2.MustNewRedis(c)}
}
