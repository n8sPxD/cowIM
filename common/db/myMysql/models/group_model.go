package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Groupname   string       `json:"groupname"`
	Members     []uint32     `json:"members"`
	GroupConfig *GroupConfig `gorm:"foreignKey:GroupID" json:"groupConfig"`
}

type GroupConfig struct {
	gorm.Model
	GroupID uint32 `json:"groupID"`
	Group   Group  `gorm:"foreignKey:GroupID" json:"-"`
}
