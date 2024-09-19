package models

import (
	"time"
)

type Message struct {
	ID       int64     `bson:"_id,omitempty"      json:"id,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
	FromID   uint32    `bson:"from"               json:"from"`
	ToID     uint32    `bson:"to"                 json:"to"`
	Content  string    `bson:"content"            json:"content"`
}
