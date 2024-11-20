package models

import (
	"gorm.io/gorm"
	"time"
)

type Friends struct {
	gorm.Model
	ID        uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uint32    `json:"userID"`
	FriendID  uint32    `json:"friendID"`
	// TODO: FriendNote string `json:"friendNote"` // 好友备注
	Friend User `gorm:"foreignKey:FriendID" json:"-"`
}
