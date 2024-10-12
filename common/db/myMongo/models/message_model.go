package models

import "time"

// MessageSync 消息同步表，用于用户即时能查询的信息，不直接入库，由Timeline间接入库（做数据冗余）
type MessageSync struct {
	ID        int64     `bson:"_id"              json:"id"`
	MsgType   uint8     `bson:"msg_type"         json:"msgType"`
	Content   string    `bson:"content"          json:"content"`
	Extend    int64     `bson:"extend,omitempty" json:"extend,omitempty"`
	Timestamp time.Time `bson:"timestamp"        json:"timestamp"` // 用于删除过时消息
}

// MessageRecord 消息记录表，用于消息的持久化
type MessageRecord struct {
	ID         int64     `bson:"_id"              json:"id"` // 此处id全局唯一
	SenderID   uint32    `bson:"sender_id"        json:"senderID"`
	Type       uint8     `bson:"type"             json:"type"` // 单聊或者群聊
	ReceiverID uint32    `bson:"receiver_id"      json:"receiverID"`
	MsgType    uint8     `bson:"msg_type"         json:"msgType"`
	Content    string    `bson:"content"          json:"content"`
	Extend     int64     `bson:"extend,omitempty" json:"extend,omitempty"`
	Timestamp  time.Time `bson:"timestamp"        json:"timestamp"` // ps:基于timestamp删除过期数据
}

/*
	单聊:
		A --哈哈--> B
		A.UserID = 233, B.UserID = 666
	示例存表:
	{
		ID: 12345678910, // 使用idgen.NextId()生成
		SenderID: 233,
		Type: 0, // 单聊的常量表示
		ReceiverID: 666,
		MsgType: 0, // 一般文本消息的常量表示
		Content: "哈哈",
		Timestamp: time.Now()
	}

   -----------------------------------------------------------

	群聊:
		GroupA = [ A, B, C ]
		A --我嘞个豆--> GroupA
		A.UserID = 233, GroupAID = 114514
	示例存表:
	{
		ID: 233333333333,  // 使用idgen.NextId()生成
		SenderID: 233,
		Type: 1, // 群聊的常量表示
		ReceiverID: 114514,
		MsgType: 0,
		Content: "我嘞个豆",
		Timestamp: time.Now()
	}

*/
