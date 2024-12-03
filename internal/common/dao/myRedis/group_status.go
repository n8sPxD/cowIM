package myRedis

import (
	"context"
	"strconv"
)

func (db *Native) GetGroupMembers(ctx context.Context, id uint32) ([]string, error) {
	if ids, err := db.SMembers(ctx, "group_"+strconv.Itoa(int(id))).Result(); err != nil {
		return []string{}, err
	} else {
		return ids, nil
	}
}

func (db *Native) AddGroupMembers(ctx context.Context, groupid uint32, userid []uint32) error {
	_, err := db.SAdd(ctx, "group_"+strconv.Itoa(int(groupid)), userid).Result()
	return err
}

func (db *Native) RemGroupMembers(ctx context.Context, groupid uint32, userid []uint32) error {
	_, err := db.SRem(ctx, strconv.Itoa(int(groupid)), userid).Result()
	return err
}
