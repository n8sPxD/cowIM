package logic

import (
	"context"
	"errors"
	"time"

	"github.com/n8sPxD/cowIM/internal/apis/info/internal/svc"
	"github.com/n8sPxD/cowIM/internal/apis/info/internal/types"
	"github.com/n8sPxD/cowIM/internal/common/db/myMongo/models"
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

type TimelineSyncResponse struct {
	Infos []models.UserTimeline `json:"infos"`
}

// TimelineSync 以用户本地最新消息为起始点，获取服务器中更新的消息
func (l *TimelineSyncLogic) TimelineSync(req *types.TimelineSyncRequest) (*TimelineSyncResponse, error) {
	// 获取消息
	chats, err := l.svcCtx.Mongo.GetRecentChatList(l.ctx, req.ID, time.Unix(req.Timestamp, 0))
	if err != nil {
		logx.Error("[TimelineSync] Get recent chat from db failed, err: ", err)
		return nil, errors.New("获取消息失败")
	}

	return &TimelineSyncResponse{Infos: chats}, nil
}
