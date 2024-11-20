package myRedis

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
)

func (db *DB) GetUserRouterStatus(ctx context.Context, userID uint32) (string, error) {
	return db.Get(ctx, fmt.Sprintf("r_%d", userID)).Result()
}

func (db *DB) UpdateUserRouterStatus(ctx context.Context, userID uint32, workID uint16) (string, error) {
	// TODO: 最后的expiration根据用户心跳时间来修改
	return db.Set(ctx, fmt.Sprintf("r_%d", userID), strconv.Itoa(int(workID)), 0).Result()
}

func (db *DB) RemoveUserRouterStatus(ctx context.Context, userID uint32) (int64, error) {
	return db.Del(ctx, fmt.Sprintf("r_%d", userID)).Result()
}

func (db *DB) RemoveAllUserRouterStatus(ctx context.Context, workerID uint16) error {
	var cursor uint64
	for {
		var (
			keys []string
			next uint64
			err  error
		)
		if keys, next, err = db.Scan(ctx, cursor, "r_*", 100).Result(); err != nil {
			return err
		}
		for _, key := range keys {
			if v, err := db.Get(ctx, key).Result(); err != nil {
				if !errors.Is(err, redis.Nil) {
					return err
				}
				continue
			} else {
				if v == strconv.Itoa(int(workerID)) {
					if _, err := db.Unlink(ctx, key).Result(); err != nil {
						return err
					}
				}
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}
