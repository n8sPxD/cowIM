package models

import (
	"time"
)

type User struct {
	ID         uint32      `gorm:"primaryKey" json:"ID"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
	Username   string      `gorm:"size:64"           json:"username"`
	Avatar     string      `                         json:"avatar"`
	Password   string      `                         json:"-"`
	UserConfig *UserConfig `gorm:"foreignKey:UserID" json:"UserConfig"`
}

type UserConfig struct {
	UserID    uint32    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	User      User      `json:"-"      gorm:"foreignKey:UserID" `
}
