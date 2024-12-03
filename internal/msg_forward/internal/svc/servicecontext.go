package svc

import (
	"fmt"
	"github.com/n8sPxD/cowIM/internal/business/group/rpc/group"
	"github.com/n8sPxD/cowIM/internal/business/group/rpc/types/groupRpc"
	"github.com/n8sPxD/cowIM/internal/common/dao/myMongo"
	"github.com/n8sPxD/cowIM/internal/common/dao/myMysql"
	"github.com/n8sPxD/cowIM/internal/common/dao/myRedis"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/config"
	"github.com/n8sPxD/cowIM/pkg/servicehub"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	Redis  *myRedis.Native
	Mongo  *myMongo.DB

	MySQL  *myMysql.DB
	Regist *servicehub.RegisterHub
	Discov *servicehub.DiscoveryHub

	GroupRpc groupRpc.GroupClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Redis:    myRedis.MustNewNativeRedis(c.RedisConf),
		Mongo:    myMongo.MustNewMongo(fmt.Sprintf("mongodb://%s", c.MongoConf.Host)),
		MySQL:    myMysql.MustNewMySQL(c.MySQL.DataSource),
		Regist:   servicehub.NewRegisterHub(c.Etcd.Endpoints, 3),
		Discov:   servicehub.NewDiscoveryHub(c.Etcd.Endpoints),
		GroupRpc: group.NewGroup(zrpc.MustNewClient(c.GroupRpc)),
	}
}
