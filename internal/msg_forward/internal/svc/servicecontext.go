package svc

import (
	"fmt"
	"github.com/n8sPxD/cowIM/internal/common/dao/myMongo"
	"github.com/n8sPxD/cowIM/internal/common/dao/myMysql"
	"github.com/n8sPxD/cowIM/internal/common/dao/myRedis"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/config"
	"github.com/n8sPxD/cowIM/pkg/servicehub"
)

type ServiceContext struct {
	Config config.Config
	Redis  *myRedis.Native
	Mongo  *myMongo.DB

	MySQL  *myMysql.DB
	Regist *servicehub.RegisterHub
	Discov *servicehub.DiscoveryHub
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  myRedis.MustNewNativeRedis(c.RedisConf),
		Mongo:  myMongo.MustNewMongo(fmt.Sprintf("mongodb://%s", c.MongoConf.Host)),
		//MsgDBSaver: &kafka.Writer{
		//	Addr:         kafka.TCP(c.MsgDBSaver.Brokers...),
		//	Topic:        c.MsgDBSaver.Topic,
		//	Balancer:     &kafka.LeastBytes{},
		//	BatchTimeout: 10 * time.Millisecond, // 低超时时间
		//	RequiredAcks: kafka.RequireOne,      // 仅等待 Leader 确认
		//	Compression:  compress.Zstd,         // Zstd压缩
		//	Async:        true,                  // 启用异步写入
		//	MaxAttempts:  1,                     // 限制重试次数
		//},
		MySQL:  myMysql.MustNewMySQL(c.MySQL.DataSource),
		Regist: servicehub.NewRegisterHub(c.Etcd.Endpoints, 3),
		Discov: servicehub.NewDiscoveryHub(c.Etcd.Endpoints),
	}
}
