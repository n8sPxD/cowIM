package logic

import (
	"context"
	"errors"
	"github.com/n8sPxD/cowIM/internal/business/group/rpc/internal/svc"
	"github.com/n8sPxD/cowIM/internal/business/group/rpc/types/groupRpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMembersLogic {
	return &GetGroupMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupMembers 从数据库获取群员数据
func (l *GetGroupMembersLogic) GetGroupMembers(in *groupRpc.GroupMembersRequest) (*groupRpc.GroupMembersResponse, error) {
	id := in.GetId()
	if id == 0 {
		return nil, errors.New("id must set")
	}

	ids, err := l.svcCtx.DB.GetGroupMemberIDs(l.ctx, id)
	if err != nil {
		return nil, err
	}

	return &groupRpc.GroupMembersResponse{Ids: ids}, nil
}
