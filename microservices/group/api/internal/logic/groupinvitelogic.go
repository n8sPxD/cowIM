package logic

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/microservices/group/api/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInviteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInviteLogic {
	return &GroupInviteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupInvite 拉人进群，暂时定为邀请人就能直接拉进群而不需要别人同意
func (l *GroupInviteLogic) GroupInvite(req *types.GroupInviteRequest) (*types.GroupInviteResponse, error) {
	// 直接塞mysql里
	err := l.svcCtx.MySQL.InsertGroupMembers(l.ctx, req.GroupID, req.Members)
	if err != nil {
		logx.Errorf("[GroupInvite] Insert members to group %d failed, error: %v", req.GroupID, err)
		return nil, errors.New("邀请成员入群失败")
	}
	return nil, nil
}
