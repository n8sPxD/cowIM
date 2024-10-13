package test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestBloom(t *testing.T) {
	store := redis.MustNewRedis(redis.RedisConf{
		Host: "localhost:6379",
		Type: redis.NodeType,
	})
	filter := bloom.New(store, "test", 100)
	test1 := uuid.NewString()
	filter.Add([]byte(test1))
	ok, _ := filter.Exists([]byte(test1))
	if !ok {
		t.Fail()
	}
	ok2, _ := filter.Exists([]byte(uuid.NewString()))
	if ok2 {
		t.Fail()
	}
}
