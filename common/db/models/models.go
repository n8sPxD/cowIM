package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string      `gorm:"unique" json:"username"`
	Password   string      `json:"password"`
	UserConfig *UserConfig `gorm:"foreignKey:UserID" json:"UserConfig"`
}

type UserConfig struct {
	gorm.Model
	UserID uint `json:"userID"`
	User   User `gorm:"foreignKey:UserID" json:"-"`
}
