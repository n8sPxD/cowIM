package myRedis

import (
	"context"
	"time"
)

func (db *DB) CheckDuplicateMessage(ctx context.Context, uuid string) (bool, error) {
	var (
		key = "dup_id"
		ok  bool
		err error
	)
	if ok, err = db.SismemberCtx(ctx, key, uuid); err != nil { // redis 出问题了
		return false, err
	} else if !ok { // 没找到该元素，说明该消息是第一次到服务器
		db.SaddCtx(ctx, key, uuid)
		go func() {
			time.Sleep(5 * time.Second)
			db.Srem(key, uuid)
		}()
	}
	return ok, nil
}

func (db *DB) RemoveAllDupMessages() {
	fields, _ := db.Smembers("dup_id")
	db.Srem("dup_id", fields)
}
