package models

import (
	"time"
)

type Group struct {
	ID           uint32        `gorm:"primaryKey" json:"groupId"`
	CreateAt     time.Time     `json:"createAt"`
	UpdateAt     time.Time     `json:"updateAt"`
	GroupName    string        `json:"groupName"`
	GroupAvatar  string        `json:"groupAvatar"`
	GroupMembers []GroupMember `gorm:"foreignKey:GroupID" json:"-"`
}

type GroupMember struct {
	ID          uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	GroupID     uint32    `json:"groupID"`
	Group       Group     `gorm:"foreignKey:GroupID" json:"-"`
	UserID      uint32    `json:"userID"`
	IngroupName string    `json:"ingroupName"`
	Role        int8      `json:"role"`
}
