package db

import (
	"context"

	"github.com/n8sPxD/cowIM/common/db/models"
)

func (db *DB) InsertUser(ctx context.Context, user *models.User) error {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	err := db.client.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	conf := models.UserConfig{
		UserID: user.ID,
	}
	return db.client.WithContext(ctx).Create(&conf).Error
}
