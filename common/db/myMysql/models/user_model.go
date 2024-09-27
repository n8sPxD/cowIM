package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string      `gorm:"size:64"           json:"username"`
	Avatar     string      `                         json:"avatar"`
	Password   string      `                         json:"password"`
	UserConfig *UserConfig `gorm:"foreignKey:UserID" json:"UserConfig"`
}

type UserConfig struct {
	gorm.Model
	UserID uint32 `json:"userID"`
	User   User   `json:"-"      gorm:"foreignKey:UserID"`
}
