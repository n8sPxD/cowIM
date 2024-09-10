package db

import (
	"context"

	"github.com/n8sPxD/cowIM/common/db/models"
)

func (db *DB) InsertUser(ctx context.Context, user *models.User) error {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	if err := db.client.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return db.client.WithContext(ctx).Create(&models.UserConfig{UserID: user.ID}).Error
}
