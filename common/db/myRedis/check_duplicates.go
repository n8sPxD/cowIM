package myRedis

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/bloom"
)

func (db *DB) CheckDuplicateMessage(ctx context.Context, uuid string) (bool, error) {
	var (
		filter *bloom.Filter
		ok     bool
		err    error
	)
	// 没有初始化该filter
	if filter, ok = db.blooms["dup_id"]; !ok {
		return false, errors.New("failed to get bloom filter \"dup_id\"")
	}
	if ok, err = filter.Exists([]byte(uuid)); err != nil {
		// 查找过程出现错误
		return false, err
	} else if !ok {
		// 没找到，不存在
		filter.Add([]byte(uuid))
	}
	return ok, nil
}
