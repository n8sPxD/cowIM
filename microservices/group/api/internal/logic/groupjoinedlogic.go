package logic

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/microservices/group/api/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupJoinedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupJoinedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupJoinedLogic {
	return &GroupJoinedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupJoinedLogic) GroupJoined(req *types.GroupJoinedRequest) (*types.GroupJoinedResponse, error) {
	// 直接查数据库
	groups, err := l.svcCtx.MySQL.GetGroupIDJoined(l.ctx, req.UserID)
	if err != nil {
		logx.Error("[GroupJoined] Get groups from MySQL failed, error: ", err)
		return nil, errors.New("获取群组信息失败")
	}
	return &types.GroupJoinedResponse{GroupID: groups}, nil
}