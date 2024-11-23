package myRedis

import (
	"context"
	"time"
)

func (db *Native) CheckDuplicateMessage(ctx context.Context, uuid string) (bool, error) {
	var (
		key = "dup_id"
		ok  bool
		err error
	)
	if ok, err = db.SIsMember(ctx, key, uuid).Result(); err != nil { // redis 出问题了
		return false, err
	} else if !ok { // 没找到该元素，说明该消息是第一次到服务器
		db.SAdd(ctx, key, uuid)
		go func() {
			time.Sleep(5 * time.Second)
			db.SRem(ctx, key, uuid)
		}()
	}
	return ok, nil
}

func (db *Native) RemoveAllDupMessages(ctx context.Context) {
	fields, _ := db.SMembers(ctx, "dup_id").Result()
	db.SRem(ctx, "dup_id", fields)
}
