package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/n8sPxD/cowIM/microservices/auth/rpc/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/auth/rpc/types/authRpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserConnStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserConnStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserConnStatusLogic {
	return &UserConnStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserConnStatusLogic) UserConnStatus(in *authRpc.ConnRequest) (*authRpc.Empty, error) {
	// 路由信息登记
	// Key: user_id 			Value: { server_work_id: xxx, last_update: xxx }
	// 用户路由信息，保存用户建立长连接的服务器IP和最后和服务器进行心跳检测的时间
	type status struct {
		WorkID     uint32    `json:"work_id"`
		LastUpdate time.Time `json:"last_update"`
	}
	tmp := status{
		WorkID:     in.WorkId,
		LastUpdate: time.Now(),
	}
	val, err := json.Marshal(tmp)
	if err != nil {
		logx.Error("[UserRouteStatus] Json marshal failed, error: ", err)
		return nil, err
	}

	key := fmt.Sprintf("%d_route", in.UserId)
	if err := l.svcCtx.Redis.HsetCtx(l.ctx, "router", key, string(val)); err != nil {
		logx.Error("[UserRouteStatus] Redis Hset failed, error: ", err)
		return nil, err
	}

	return &authRpc.Empty{}, nil
}
