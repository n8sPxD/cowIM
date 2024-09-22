package myRedis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type Status struct {
	WorkID     uint32    `json:"work_id"`
	LastUpdate time.Time `json:"last_update"`
}

func (db *DB) GetUserRouterStatus(ctx context.Context, userID uint32) (*Status, error) {
	logx.Debug("router key: ", strconv.FormatInt(int64(userID), 10))
	res, err := db.HgetCtx(ctx, "router", strconv.FormatInt(int64(userID), 10))
	if err != nil {
		return nil, err
	}
	var status Status
	err = json.Unmarshal([]byte(res), &status)
	if err != nil {
		return nil, err
	}
	return &status, err
}

func (db *DB) RemoveUserRouterStatus(ctx context.Context, userID uint32) error {
	logx.Debug("router key: ", strconv.FormatInt(int64(userID), 10))
	_, err := db.HdelCtx(ctx, "router", strconv.FormatInt(int64(userID), 10))
	if err != nil {
		return err
	}
	return nil
}