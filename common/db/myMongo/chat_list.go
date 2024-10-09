package myMongo

import (
	"context"
	"fmt"
	"time"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/db/myMongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatListInfo struct {
	SenderID  uint32 `json:"senderID"`
	GroupID   uint32 `json:"groupID,omitempty"`
	RecentMsg string `json:"recentMsg"`
}

func (db *DB) GetRecentChatList(ctx context.Context, id uint32, latest time.Time) ([]ChatListInfo, error) {
	// 初始化聚合管道
	filter := mongo.Pipeline{}

	// 首先，根据用户 ID 过滤文档
	matchStage := bson.D{
		{"$match", bson.D{
			{"receiver_id", id},
		}},
	}
	filter = append(filter, matchStage)

	// 如果 latest 时间戳不为空，则添加 timestamp 的过滤条件
	if latest != time.Unix(0, 0) {
		timestampMatch := bson.D{
			{"$match", bson.D{
				{"timestamp", bson.D{{"$gte", latest}}},
			}},
		}
		filter = append(filter, timestampMatch)
	}

	// 按照 timestamp 倒序排序
	sortStage := bson.D{
		{"$sort", bson.D{
			{"timestamp", -1},
		}},
	}
	filter = append(filter, sortStage)

	// 按照 id 分组，并获取每组的最新记录
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$id"}, // 按照 id 字段分组
			{"latestRecord", bson.D{{"$first", "$$ROOT"}}}, // 获取每组最新的记录
		}},
	}
	filter = append(filter, groupStage)

	// 将嵌套文档替换为根文档
	replaceRootStage := bson.D{
		{"$replaceRoot", bson.D{
			{"newRoot", "$latestRecord"},
		}},
	}
	filter = append(filter, replaceRootStage)

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
		var chat ChatListInfo
		chat.SenderID = timeline.SenderID
		chat.RecentMsg = getMsgPreview(timeline)
		chat.GroupID = timeline.GroupID // GroupID默认值为0
		chatList = append(chatList, chat)
	}

	return chatList, nil
}

func getMsgPreview(msg models.UserTimeline) string {
	content := msg.Message.Content
	switch msg.Message.MsgType {
	case constant.MSG_COMMON_MSG:
		runes := []rune(content)
		if len(runes) <= 50 {
			return content
		}
		return string(runes[:50]) + "..."
	default:
		return ""
	}
}
