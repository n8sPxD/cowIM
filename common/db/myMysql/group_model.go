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
