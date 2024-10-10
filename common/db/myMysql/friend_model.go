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
	return db.client.WithContext(ctx).Create(&relation).Error
}

func (db *DB) GetFriends(ctx context.Context, userID uint32) ([]models.User, error) {
	var friends []models.User
	// 子查询，分别对应自己加别人 和 别人加自己的情况
	sub1 := db.client.
		Table("friends").
		Select("friend_id").
		Where("user_id = ?", userID)
	sub2 := db.client.
		Table("friends").
		Select("user_id").
		Where("friend_id = ?", userID)

	err := db.client.
		Debug().
		Select("username", "id").
		Where("id = (?)", sub1).
		Or("id = (?)", sub2).
		Find(&friends).
		Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}
