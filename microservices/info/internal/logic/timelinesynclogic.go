package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/microservices/info/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/info/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TimelineSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTimelineSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TimelineSyncLogic {
	return &TimelineSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// TimelineSync 以用户本地最新消息为起始点，获取服务器中更新的消息
func (l *TimelineSyncLogic) TimelineSync(req *types.TimelineSyncRequest) (resp *types.TimelineSyncResponse, err error) {
	messages, err :=
	return
}
