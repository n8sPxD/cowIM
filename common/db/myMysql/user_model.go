package myMysql

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/common/db/myMysql/models"
	"gorm.io/gorm"
)

// InsertUser 创建新的用户，并判断用户名重复
func (db *DB) InsertUser(ctx context.Context, user *models.User) error {
	if err := db.client.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}
		return err
	}
	return db.client.WithContext(ctx).Create(&models.UserConfig{UserID: uint(user.ID)}).Error
}

// GetUserInfo 获取特定用户的信息（密码、用户名）
func (db *DB) GetUserInfo(ctx context.Context, id uint32) (*models.User, error) {
	var user models.User
	err := db.client.WithContext(ctx).
		Model(&models.User{}).
		Select("password", "username").
		Where("id = ?", id).
		Take(&user).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserPassword 获取特定用户的密码
func (db *DB) GetUserPassword(ctx context.Context, id uint32) (*string, error) {
	var password string
	err := db.client.WithContext(ctx).
		Model(&models.User{}).
		Select("password").
		Where("id = ?", id).
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
func (db *DB) GetUser(ctx context.Context, id uint32) (*models.User, error) {
	var user models.User
	err := db.client.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
