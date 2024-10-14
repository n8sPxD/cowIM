package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/internal/apis/wsget/api/internal/svc"
	"github.com/n8sPxD/cowIM/internal/apis/wsget/api/internal/types"
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

func (l *GetAvailableWSServerLogic) GetAvailableWSServer(req *types.WebsocketServerGetRequest) (*types.WebsocketServerGetResponse, error) {
	// 直接从etcd获取服务器ip返回
	return &types.WebsocketServerGetResponse{IP: l.svcCtx.DiscoveryHub.GetServiceEndpoint(l.ctx)}, nil
}
