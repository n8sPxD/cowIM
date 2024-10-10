package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/microservices/wsget/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/wsget/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvailableWSServerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvailableWSServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvailableWSServerLogic {
	return &GetAvailableWSServerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvailableWSServerLogic) GetAvailableWSServer(req *types.WebsocketServerGetRequest) (resp *types.WebsocketServerGetResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
