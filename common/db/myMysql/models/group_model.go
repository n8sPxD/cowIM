package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	GroupName    string       `json:"groupName"`
	GroupMembers []User       `gorm:"many2many:group_user" json:"groupMembers"`
	GroupConfig  *GroupConfig `json:"groupConfig"`
}

type GroupConfig struct {
	gorm.Model
	GroupID uint32 `json:"groupID"`
	Group   Group  `gorm:"foreignKey:GroupID" json:"-"`
}
