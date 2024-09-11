package db

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/common/db/models"
	"gorm.io/gorm"
)

// InsertUser 创建新的用户，并判断用户名重复
func (db *DB) InsertUser(ctx context.Context, user *models.User) error {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	if err := db.client.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}
		return err
	}
	return db.client.WithContext(ctx).Create(&models.UserConfig{UserID: user.ID}).Error
}

// GetUserPassword 获取特定用户的密码
func (db *DB) GetUserPassword(ctx context.Context, username string) (*string, error) {
	db.rwMutex.RLock()
	defer db.rwMutex.RUnlock()
	var password string
	err := db.client.WithContext(ctx).
		Model(&models.User{}).
		Select("password").
		Where("username = ?", username).
		Take(&password).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &password, nil
}

// GetUser 获取特定用户所有信息
func (db *DB) GetUser(ctx context.Context, username string) (*models.User, error) {
	db.rwMutex.RLock()
	defer db.rwMutex.RUnlock()
	var user models.User
	err := db.client.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
