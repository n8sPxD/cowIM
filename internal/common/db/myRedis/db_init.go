package myRedis

import (
	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type DB struct {
	*redis.Redis
	blooms map[string]*bloom.Filter
}

func MustNewRedis(c redis.RedisConf) *DB {
	return &DB{
		Redis:  redis.MustNewRedis(c),
		blooms: make(map[string]*bloom.Filter),
	}
}

func (db *DB) WithBloom(key string, bits uint) *DB {
	db.blooms[key] = bloom.New(db.Redis, key, bits)
	return db
}
