package svc

import (
	"fmt"
	"time"

	"github.com/n8sPxD/cowIM/common/db/myMongo"
	"github.com/n8sPxD/cowIM/common/db/myRedis"
	"github.com/n8sPxD/cowIM/microservices/msgForward/internal/config"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type ServiceContext struct {
	Config     config.Config
	Redis      *myRedis.DB
	Mongo      *myMongo.DB
	MsgSender  *kafka.Writer
	MsgDBSaver *kafka.Writer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  myRedis.MustNewRedis(c.RedisConf),
		Mongo:  myMongo.MustNewMongo(fmt.Sprintf("mongodb://%s", c.MongoConf.Host)),
		MsgSender: &kafka.Writer{
			Addr:         kafka.TCP(c.MsgSender.Brokers...),
			Topic:        c.MsgSender.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond, // 低超时时间
			RequiredAcks: kafka.RequireOne,      // 仅等待 Leader 确认
			Compression:  compress.Zstd,         // Zstd压缩
			Async:        true,                  // 启用异步写入
			MaxAttempts:  1,                     // 限制重试次数
		},
		MsgDBSaver: &kafka.Writer{
			Addr:         kafka.TCP(c.MsgDBSaver.Brokers...),
			Topic:        c.MsgDBSaver.Topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond, // 低超时时间
			RequiredAcks: kafka.RequireOne,      // 仅等待 Leader 确认
			Compression:  compress.Zstd,         // Zstd压缩
			Async:        true,                  // 启用异步写入
			MaxAttempts:  1,                     // 限制重试次数
		},
	}
}
