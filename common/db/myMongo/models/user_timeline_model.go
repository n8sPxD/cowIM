package models

import "time"

// UserTimeline 用户时间线
type UserTimeline struct {
	ID         int64       `bson:"_id"                json:"id"`
	ReceiverID uint32      `bson:"receiver_id"        json:"receiverID"`
	SenderID   uint32      `bson:"sender_id"          json:"senderID"`
	GroupID    uint32      `bson:"group_id,omitempty" json:"groupID,omitempty"`
	Message    MessageSync `bson:"msgForward"         json:"msgForward"`
	Timestamp  time.Time   `bson:"timestamp"          json:"timestamp"` // 用于删除过时消息 + 实现Timeline模型(用户消息按时间线排列)
}

/*
	单聊:
		A --我去--> B
		A.UserID = 233, B.UserID = 666
		UserID = 233, SenderID = 666
	示例存表, 存B的timeline:
	{
		ID: 11451419198, // 使用idgen.NextId()生成
		UserID: 666,
		SenderID: 233,
		Message: MessageSync,
		Timestamp: time.Now()
	}
	MessageSync:
	{
		ID: 1,
		MsgType: 0, // 一般文本消息的常量表示
		Content: "我去"
		Timestamp: time.Now()
	}

   ----------------------------------------------------

	群聊:
		GroupA = [ A, B, C ]
		A --无语--> GroupA
		A.UserID = 233, B.UserID = 666, C.UserID = 213, GroupAID = 114514
	示例存表, 存B和C的timeline:
	B:
	{
		ID: 12345678987,
		UserID: 666,
		SenderID: 233,
		GroupID: 114514,
		Message: MessageSync
		Timestamp: time.Now()
	}
	C:
	{
		ID: 12345678988,
		UserID: 213,
		SenderID: 233,
		GroupID: 114514,
		Message: MessageSync
		Timestamp: time.Now()
	}
	MessageSync:
	{
		ID: 2,
		MsgType: 0,
		Content: "无语",
		Timestamp: time.Now()
	}
*/
