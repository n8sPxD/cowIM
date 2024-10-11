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

func (db *DB) GetRecentChatList(ctx context.Context, id uint32, latest time.Time) ([]models.UserTimeline, error) {
	// 初始化聚合管道
	filter := mongo.Pipeline{}

	// 构造 $or 逻辑，匹配 receiver_id 或 sender_id
	matchStage := bson.D{
		{"$match", bson.D{
			{"$or", bson.A{
				bson.D{{"receiver_id", id}},
				bson.D{{"sender_id", id}},
			}},
		}},
	}
	filter = append(filter, matchStage)

	// 如果 latest 时间戳不为空，则添加 timestamp 的过滤条件
	if latest != time.Unix(0, 0) {
		timestampMatch := bson.D{
			{"$match", bson.D{
				{"timestamp", bson.D{{"$gt", latest}}},
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

	// 定义结果 slice
	var timelines []models.UserTimeline

	// 执行查询
	err := db.TimeLine.Aggregate(ctx, &timelines, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user timeline: %w", err)
	}

	return timelines, nil
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
