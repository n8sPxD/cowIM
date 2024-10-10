package models

import "gorm.io/gorm"

type Friends struct {
	gorm.Model
	UserID   uint32 `json:"userID"`
	FriendID uint32 `json:"friendID"`
	Friend   User   `gorm:"foreignKey:FriendID" json:"-"`
}
