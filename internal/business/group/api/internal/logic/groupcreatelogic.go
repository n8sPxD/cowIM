package logic

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/internal/business/group/api/internal/svc"
	"github.com/n8sPxD/cowIM/internal/business/group/api/internal/types"
	"github.com/n8sPxD/cowIM/internal/common/dao/myMysql/models"
	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupCreateLogic) GroupCreate(req *types.GroupCreateRequest) (resp *types.GroupCreateResponse, err error) {
	group := models.Group{
		GroupName: req.Groupname,
	}
	if err := l.svcCtx.MySQL.InsertGroup(l.ctx, &group); err != nil {
		logx.Error("[Create] Insert group to DB failed, error: ", err)
		return nil, errors.New("创建群组失败!可能是服务器出了问题")
	}
	return &types.GroupCreateResponse{GroupID: uint32(group.ID)}, nil
}
