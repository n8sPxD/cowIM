package myMysql

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/internal/common/constant"
	"github.com/n8sPxD/cowIM/internal/common/db/myMysql/models"
	"gorm.io/gorm"
)

func (db *DB) InsertGroup(ctx context.Context, group *models.Group) error {
	if err := db.client.WithContext(ctx).Create(group).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}
		return err
	}
	return nil
}

//func (db *DB) GetGroupMembers(ctx context.Context, id uint) ([]models.GroupMember, error) {
//	var members []models.GroupMember
//	err := db.client.WithContext(ctx).
//		Model(&models.Group{}).
//		Preload("GroupMembers").
//		Select("group_members").
//		Where("id = ?", id).
//		Find(&members).
//		Error
//	return members, err
//}

func (db *DB) GetGroupMemberIDs(ctx context.Context, id uint) ([]uint, error) {
	var ids []uint
	err := db.client.WithContext(ctx).
		Model(&models.GroupMember{}).
		Select("user_id").
		Where("group_id = ?", id).
		Find(&ids).
		Error
	return ids, err
}

func (db *DB) InsertGroupMember(ctx context.Context, groupID uint32, member uint32) error {
	membercol := models.GroupMember{
		GroupID: uint32(groupID),
		UserID:  uint32(member),
		Role:    constant.GROUP_COMMON,
	}
	return db.client.WithContext(ctx).Create(&membercol).Error
}

func (db *DB) InsertGroupMembers(ctx context.Context, groupID uint32, members []uint32) error {
	membercols := make([]models.GroupMember, 0, len(members))
	for i := range membercols {
		tmp := models.GroupMember{
			GroupID: uint32(groupID),
			UserID:  uint32(members[i]),
			Role:    constant.GROUP_COMMON,
		}
		membercols[i] = tmp
	}
	return db.client.WithContext(ctx).Create(&membercols).Error
}

func (db *DB) GetGroupsJoinedBaseInfo(ctx context.Context, userID uint32) ([]models.Group, error) {
	var groups []models.Group
	err := db.client.
		WithContext(ctx).
		Model(&models.Group{}).
		Select("groups.id", "groups.group_name", "groups.group_avatar").
		Joins("JOIN group_members on groups.id = group_members.group_id").
		Where("group_members.user_id = ?", userID).
		Scan(&groups).
		Error
	return groups, err
}
