package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	GroupName    string         `json:"groupName"`
	GroupAvatar  string         `json:"groupAvatar"`
	GroupMembers []*GroupMember `json:"groupMembers"`
	GroupConfig  *GroupConfig   `json:"groupConfig"`
}

type GroupConfig struct {
	gorm.Model
	GroupID uint  `json:"groupID"`
	Group   Group `gorm:"foreignKey:GroupID" json:"-"`
}

type GroupMember struct {
	gorm.Model
	GroupID     uint   `json:"groupID"`
	Group       Group  `gorm:"foreignKey:GroupID" json:"-"`
	UserID      uint   `json:"userID"`
	IngroupName string `json:"ingroupName"`
	Role        int8   `json:"role"`
}
