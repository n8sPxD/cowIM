package myRedis

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
)

func (db *DB) GetUserRouterStatus(ctx context.Context, userID uint32) (string, error) {
	return db.GetCtx(ctx, fmt.Sprintf("r_%d", userID))
}

func (db *DB) UpdateUserRouterStatus(ctx context.Context, userID uint32, workID uint16) error {
	return db.SetCtx(ctx, fmt.Sprintf("r_%d", userID), strconv.Itoa(int(workID)))
}

func (db *DB) RemoveUserRouterStatus(ctx context.Context, userID uint32) (int, error) {
	return db.DelCtx(ctx, fmt.Sprintf("r_%d", userID))
}

func (db *DB) RemoveAllUserRouterStatus(ctx context.Context, workerID uint16) error {
	var cursor uint64
	for {
		var (
			keys []string
			next uint64
			err  error
		)
		if keys, next, err = db.ScanCtx(ctx, cursor, "r_*", 100); err != nil {
			return err
		}
		for _, key := range keys {
			if v, err := db.GetCtx(ctx, key); err != nil {
				if !errors.Is(err, redis.Nil) {
					return err
				}
				continue
			} else {
				if v == strconv.Itoa(int(workerID)) {
					if _, err := db.UnlinkCtx(ctx, key); err != nil {
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
