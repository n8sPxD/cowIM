package myRedis

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type DB struct {
	*redis.Redis
}

func MustNewRedis(c redis.RedisConf) *DB {
	return &DB{
		Redis: redis.MustNewRedis(c),
	}
}
