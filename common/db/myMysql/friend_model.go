package myMysql

import (
	"context"

	"github.com/n8sPxD/cowIM/common/db/myMysql/models"
)

func (db *DB) InsertFriend(ctx context.Context, userID, friendID uint32) error {
	relation := models.Friends{
		UserID:   userID,
		FriendID: friendID,
	}
	return db.client.WithContext(ctx).Create(relation).Error
}

func (db *DB) GetFriends(ctx context.Context, userID uint32) ([]models.User, error) {
	var friends []models.User
	err := db.client.WithContext(ctx).
		Preload("Friend").
		Model(&models.Friends{}).
		Select("friends").
		Where("user_id", userID).
		Take(&friends).
		Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}
