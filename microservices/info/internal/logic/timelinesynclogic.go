package logic

import (
	"context"
	"errors"
	"strconv"
	"time"

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
	id, _ := strconv.Atoi(req.ID)
	// 获取消息
	chats, err := l.svcCtx.Mongo.GetRecentChatList(l.ctx, uint32(id), time.Unix(req.Timestamp, 0))
	if err != nil {
		logx.Error("[TimelineSync] Get recent chat from db failed, err: ", err)
		return nil, errors.New("获取消息失败")
	}

	// 将消息封装到resp中
	infos := make([]types.TimelineSyncInfo, len(chats))
	for i, chat := range chats {
		var info types.TimelineSyncInfo
		info.Message = chat.RecentMsg
		info.SenderID = chat.SenderID
		info.GroupID = chat.GroupID
		infos[i] = info
	}

	resp = new(types.TimelineSyncResponse)
	resp.Infos = infos

	return
}
