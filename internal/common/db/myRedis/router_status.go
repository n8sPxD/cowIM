package myRedis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
)

type RouterStatus struct {
	WorkID     uint16    `json:"work_id"`
	LastUpdate time.Time `json:"last_update"`
}

func (db *DB) GetUserRouterStatus(ctx context.Context, userID uint32) (*RouterStatus, error) {
	res, err := db.HgetCtx(ctx, "router", strconv.FormatInt(int64(userID), 10))
	if err != nil {
		return nil, err
	}
	var status RouterStatus
	err = json.Unmarshal([]byte(res), &status)
	if err != nil {
		return nil, err
	}
	return &status, err
}

// UpdateUserRouterStatus 路由信息登记
// Key: user_id			Value: { server_work_id: xxx, last_update: xxx }
// 用户路由信息，保存用户建立长连接的服务器IP和最后和服务器进行心跳检测的时间
func (db *DB) UpdateUserRouterStatus(ctx context.Context, userID uint32, workID uint16, timestamp time.Time) error {
	if val, err := json.Marshal(RouterStatus{
		WorkID:     workID,
		LastUpdate: timestamp,
	}); err != nil {
		return err
	} else {
		return db.HsetCtx(ctx, "router", strconv.FormatInt(int64(userID), 10), string(val))
	}
}

func (db *DB) RemoveUserRouterStatus(ctx context.Context, userID uint32) error {
	_, err := db.HdelCtx(ctx, "router", strconv.FormatInt(int64(userID), 10))
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RemoveAllUserRouterStatus() {
	table, _ := db.Hgetall("router")
	slices := make([]string, 0, len(table))
	for key := range table {
		slices = append(slices, key)
	}
	db.Hdel("router", slices...)
}
