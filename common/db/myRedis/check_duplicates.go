package myRedis

import (
	"context"
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
		_, _ = db.SaddCtx(ctx, key, uuid)
	}
	return ok, nil
}
