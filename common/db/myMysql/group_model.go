package myMysql

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/common/constant"
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
	return db.client.WithContext(ctx).Create(&models.GroupConfig{GroupID: group.ID}).Error
}

func (db *DB) GetGroupMembers(ctx context.Context, id uint) ([]models.GroupMember, error) {
	var members []models.GroupMember
	err := db.client.WithContext(ctx).
		Model(&models.Group{}).
		Preload("GroupMembers").
		Select("group_members").
		Where("id = ?", id).
		Find(&members).
		Error
	return members, err
}

func (db *DB) InsertGroupMember(ctx context.Context, groupID uint32, member uint32) error {
	membercol := models.GroupMember{
		GroupID: uint(groupID),
		UserID:  uint(member),
		Role:    constant.GROUP_COMMON,
	}
	return db.client.WithContext(ctx).Create(&membercol).Error
}

func (db *DB) InsertGroupMembers(ctx context.Context, groupID uint32, members []uint32) error {
	membercols := make([]models.GroupMember, 0, len(members))
	for i := range membercols {
		tmp := models.GroupMember{
			GroupID: uint(groupID),
			UserID:  uint(members[i]),
			Role:    constant.GROUP_COMMON,
		}
		membercols[i] = tmp
	}
	return db.client.WithContext(ctx).Create(&membercols).Error
}

func (db *DB) GetGroupIDJoined(ctx context.Context, userID uint32) ([]uint32, error) {
	var groups []uint32
	err := db.client.
		WithContext(ctx).
		Model(&models.GroupMember{}).
		Select("group_id").
		Where("user_id = ?", userID).
		Take(&groups).
		Error
	return groups, err
}
