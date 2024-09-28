package logic

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/common/db/myMysql/models"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.GroupCreateRequest) ( *types.GroupCreateResponse, error) {
    group := models.Group {
        Groupname: req.Groupname,
    }
    if err := l.svcCtx.MySQL.InsertGroup(l.ctx, &group); err != nil {
        logx.Error("[Create] Insert group to DB failed, error: ", err)
        return nil, errors.New("创建群组失败!可能是服务器出了问题")
    }
    return &types.GroupCreateResponse{GroupID: uint32(group.ID)}, nil
}
