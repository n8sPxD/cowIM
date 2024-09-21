package myMongo

import (
	"github.com/zeromicro/go-zero/core/stores/mon"
)

type DB struct {
	MessageRecord  *mon.Model
	TimeLine       *mon.Model
	MessageContent *mon.Model
	ChatRoom       *mon.Model
	ChatRoomUser   *mon.Model
}

func MustNewMongo(url string, opts ...mon.Option) *DB {
	return &DB{
		MessageRecord:  newModelOrPanic(url, "im_server_db", "messages_records", opts...),
		TimeLine:       newModelOrPanic(url, "im_server_db", "timeline", opts...),
		MessageContent: newModelOrPanic(url, "im_server_db", "messages", opts...),
		ChatRoom:       newModelOrPanic(url, "im_server_db", "chat_room", opts...),
		ChatRoomUser:   newModelOrPanic(url, "im_server_db", "room_user", opts...),
	}
}

func newModelOrPanic(url, dbName, collection string, opts ...mon.Option) *mon.Model {
	model, err := mon.NewModel(url, dbName, collection, opts...)
	if err != nil {
		panic("Failed to connect to MongoDB collection " + collection + ": " + err.Error())
	}
	return model
}
