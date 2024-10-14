package models

import "gorm.io/gorm"

type Friends struct {
	gorm.Model
	UserID   uint32 `json:"userID"`
	FriendID uint32 `json:"friendID"`
	// TODO: FriendNote string `json:"friendNote"` // 好友备注
	Friend User `gorm:"foreignKey:FriendID" json:"-"`
}
