package myMysql

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/common/db/myMysql/models"
	"gorm.io/gorm"
)

func (db *DB) InsertGroup(ctx context.Context, group *models.Group) error {
	if err := db.client.WithContext(ctx).Create(group).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}
		return err
	}
	return db.client.WithContext(ctx).Create(&models.GroupConfig{GroupID: uint32(group.ID)}).Error
}

func (db *DB) GetGroupMembers(ctx context.Context, id uint) ([]models.User, error) {
	var members []models.User
	err := db.client.WithContext(ctx).
		Model(&models.Group{}).
		Preload("GroupMembers").
		Select("group_members").
		Where("id = ?", id).
		Find(&members).
		Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
