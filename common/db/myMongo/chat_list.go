package myMongo

import (
	"context"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/db/myMongo/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ChatListInfo struct {
	SenderID  uint32 `json:"senderID"`
	RecentMsg string `json:"recentMsg"`
}

func (db *DB) GetRecentChatList(ctx context.Context, id uint32, maxLength int) ([]ChatListInfo, error) {
	filter := bson.M{
		"user_id": bson.M{"$eq": id},
	}
	var tl []models.UserTimeline
	if err := db.TimeLine.Find(ctx, &tl, filter); err != nil {
		return nil, err
	}
	var cl []ChatListInfo
	for _, v := range tl {
		var list ChatListInfo
		list.SenderID = v.SenderID
		list.RecentMsg = getMsgPreview(v, maxLength)
		cl = append(cl, list)
	}
	return cl, nil
}

func getMsgPreview(msg models.UserTimeline, maxLength int) string {
	content := msg.Message.Content
	switch msg.Message.MsgType {
	case constant.MSG_COMMON_MSG:
		if len(content) <= maxLength {
			return content
		}
		return content[:maxLength]
	default:
		return ""
	}
}
