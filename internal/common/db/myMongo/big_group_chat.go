package myMongo

import (
	"context"
	"fmt"
	"time"

	"github.com/n8sPxD/cowIM/internal/common/constant"
	"github.com/n8sPxD/cowIM/internal/common/db/myMongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) GetBigGroupMessage(ctx context.Context, groupID uint32, clientTimestamp time.Time) ([]models.UserTimeline, error) {
	filter := mongo.Pipeline{}

	matchStage := bson.D{
		{"$match", bson.D{
			{"$and", bson.A{
				bson.D{{"receiver_id", groupID}},
				bson.D{{"type", constant.BIG_GROUP_CHAT}},
				bson.D{{"group_id", groupID}},
			}},
		}},
	}
	filter = append(filter, matchStage)

	if clientTimestamp != time.Unix(0, 0) {
		timestampMatch := bson.D{
			{"$match", bson.D{
				{"timestamp", bson.D{{"$gt", clientTimestamp}}},
			}},
		}
		filter = append(filter, timestampMatch)
	}

	sortStage := bson.D{
		{"$sort", bson.D{
			{"timestamp", -1},
		}},
	}
	filter = append(filter, sortStage)

	var timelines []models.UserTimeline
	if err := db.TimeLine.Aggregate(ctx, &timelines, filter); err != nil {
		return nil, fmt.Errorf("failed to retrieve user timeline: %w", err)
	}
	return timelines, nil
}
