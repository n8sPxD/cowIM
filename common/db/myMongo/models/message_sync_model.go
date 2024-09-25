package models

import "time"

// MessageSync 消息同步表，用于用户即时能查询的信息
type MessageSync struct {
	ID        int64     `bson:"_id"              json:"id"`
	MsgType   uint8     `bson:"msg_type"         json:"msgType"`
	Content   string    `bson:"content"          json:"content"`
	Extend    int64     `bson:"extend,omitempty" json:"extend,omitempty"`
	Timestamp time.Time `bson:"timestamp"        json:"timestamp"` // 用于删除过时消息
}
