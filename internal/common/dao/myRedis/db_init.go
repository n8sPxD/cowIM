package myRedis

import (
	"github.com/redis/go-redis/v9"
	redis2 "github.com/zeromicro/go-zero/core/stores/redis"
	"strings"
)

type DB struct {
	redis.UniversalClient
}

func MustNewRedis(c redis2.RedisConf) *DB {
	return &DB{
		UniversalClient: redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:    strings.Split(c.Host, ","),
			Password: c.Pass,
		}),
	}
}
