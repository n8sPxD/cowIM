package myMysql

import (
	"context"

	"github.com/n8sPxD/cowIM/common/db/myMysql/models"
)

// InsertFriend 添加好友
func (db *DB) InsertFriend(ctx context.Context, userID, friendID uint32) error {
	relation := models.Friends{
		UserID:   userID,
		FriendID: friendID,
	}
	return db.client.WithContext(ctx).Create(&relation).Error
}

// GetFriend 鉴定好友关系
func (db *DB) GetFriend(ctx context.Context, userID, friendID uint32) (bool, error) {
	var exist int
	err := db.client.
		Model(&models.Friends{}).
		Select("count(*)").
		Where("user_id = ? and friend_id = ?", userID, friendID).
		Or("user_id = ? and friend_id = ?", friendID, userID).
		Take(&exist).
		Error
	if err != nil {
		return false, err
	}

	if exist == 1 {
		return true, nil
	}
	return false, nil
}

// GetFriends 获取好友列表
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
