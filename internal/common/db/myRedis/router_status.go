package myRedis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
)

type Status struct {
	WorkID     uint32    `json:"work_id"`
	LastUpdate time.Time `json:"last_update"`
}

func (db *DB) GetUserRouterStatus(ctx context.Context, userID uint32) (*Status, error) {
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
	_, err := db.HdelCtx(ctx, "router", strconv.FormatInt(int64(userID), 10))
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RemoveAllUserRouterStatus() {
	db.Hdel("router")
}
