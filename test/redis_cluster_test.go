package test

import (
	"context"
	"github.com/redis/go-redis/v9"
	zerordb "github.com/zeromicro/go-zero/core/stores/redis"
	"testing"
)

func TestGoZeroRedisCluster(t *testing.T) {
	rdb, err := zerordb.NewRedis(zerordb.RedisConf{
		Host: "127.0.0.1:6379",
		Type: zerordb.ClusterType,
	}, func(r *zerordb.Redis) {
		r.Addr = "127.0.0.1:6380"
	}, func(r *zerordb.Redis) {
		r.Addr = "127.0.0.1:6381"
	}, func(r *zerordb.Redis) {
		r.Addr = "127.0.0.1:6382"
	}, func(r *zerordb.Redis) {
		r.Addr = "127.0.0.1:6383"
	}, func(r *zerordb.Redis) {
		r.Addr = "127.0.0.1:6384"
	})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !rdb.Ping() {
		t.Error("redis ping failed")
		t.Fail()
	}
	err = rdb.Set("user666666", "password666")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestRedisCluster(t *testing.T) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"127.0.0.1:6379",
			"127.0.0.1:6380",
			"127.0.0.1:6381",
			"127.0.0.1:6382",
			"127.0.0.1:6383",
			"127.0.0.1:6384",
		},
		Password: "",
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
