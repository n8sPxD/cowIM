package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint32      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
	Username   string      `json:"username"`
	Password   string      `json:"password"`
	UserConfig *UserConfig `gorm:"foreignKey:UserID" json:"UserConfig"`
}

type UserConfig struct {
	gorm.Model
	UserID uint `json:"userID"`
	User   User `gorm:"foreignKey:UserID" json:"-"`
}
