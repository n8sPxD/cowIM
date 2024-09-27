package myMongo

import (
	"context"
	"fmt"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/db/myMongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatListInfo struct {
	SenderID  uint32 `json:"senderID"`
	RecentMsg string `json:"recentMsg"`
}

func (db *DB) GetRecentChatList(ctx context.Context, id uint32, maxLength int) ([]ChatListInfo, error) {
	filter := mongo.Pipeline{
		bson.D{{"$sort", bson.D{{"timestamp", -1}}}}, // 按照 created_at 字段倒序排序
		bson.D{{"$group", bson.D{
			{"_id", "$id"}, // 按照 id 字段分组
			{"latestRecord", bson.D{{"$first", "$$ROOT"}}}, // 获取每组最新的记录
		}}},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$latestRecord"}}}}, // 将嵌套文档替换为根文档
	}

	// 定义结果 slice
	var timelines []models.UserTimeline

	// 执行查询
	err := db.TimeLine.Aggregate(ctx, &timelines, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user timeline: %w", err)
	}

	// 创建 ChatListInfo slice，预分配空间
	chatList := make([]ChatListInfo, 0, len(timelines))

	// 遍历查询结果，处理消息
	for _, timeline := range timelines {
		chatList = append(chatList, ChatListInfo{
			SenderID:  timeline.SenderID,
			RecentMsg: getMsgPreview(timeline, maxLength),
		})
	}

	return chatList, nil
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
