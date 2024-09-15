package models

import "time"

type SingleMessage struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	SendUser  string    `json:"sendUser"`
	RecvUser  string    `json:"recvUser"`
	Type      int8      `json:"type"`
	Content   string    `json:"content"`
}

type GroupMessage struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	SendUser  string    `json:"sendUser"`
	RecvGroup string    `json:"recvGroup"`
	Content   string    `json:"content"`
}
