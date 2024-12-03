package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/internal/business/group/rpc/internal/svc"
	"github.com/n8sPxD/cowIM/internal/business/group/rpc/types/group"

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

func (l *GetGroupMembersLogic) GetGroupMembers(in *group.GroupMembersRequest) (*group.GroupMembersResponse, error) {
	// todo: add your logic here and delete this line

	return &group.GroupMembersResponse{}, nil
}
