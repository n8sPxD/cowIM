package svc

import (
	"fmt"

	"github.com/n8sPxD/cowIM/common/db/myMongo"
	"github.com/n8sPxD/cowIM/common/db/myRedis"
	"github.com/n8sPxD/cowIM/microservices/message/internal/config"
	"github.com/segmentio/kafka-go"
)

type ServiceContext struct {
	Config    config.Config
	Redis     *myRedis.DB
	Mongo     *myMongo.DB
	MsgSender *kafka.Writer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  myRedis.MustNewRedis(c.RedisConf),
		Mongo:  myMongo.MustNewMongo(fmt.Sprintf("mongodb://%s", c.MongoConf.Host)),
		MsgSender: &kafka.Writer{
			Addr:  kafka.TCP(c.MsgSender.Brokers...),
			Topic: c.MsgSender.Topic,
		},
	}
}
