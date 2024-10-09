package logic

import (
	"context"
	"errors"
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

type TimelineSyncInfo struct {
	SenderID  uint32    `json:"senderId"`
	GroupID   uint32    `json:"groupId,omitempty"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type TimelineSyncResponse struct {
	Infos []TimelineSyncInfo `json:"infos"`
}

// TimelineSync 以用户本地最新消息为起始点，获取服务器中更新的消息
func (l *TimelineSyncLogic) TimelineSync(req *types.TimelineSyncRequest) (resp *TimelineSyncResponse, err error) {
	// 获取消息
	chats, err := l.svcCtx.Mongo.GetRecentChatList(l.ctx, req.ID, time.Unix(req.Timestamp, 0))
	if err != nil {
		logx.Error("[TimelineSync] Get recent chat from db failed, err: ", err)
		return nil, errors.New("获取消息失败")
	}

	logx.Debug("[TimelineSync] chats: ", chats)

	// 将消息封装到resp中
	infos := make([]TimelineSyncInfo, len(chats))
	for i, chat := range chats {
		var info TimelineSyncInfo
		info.Message = chat.RecentMsg
		info.SenderID = chat.SenderID
		info.GroupID = chat.GroupID
		info.Timestamp = chat.Timestamp
		infos[i] = info
	}

	resp = new(TimelineSyncResponse)
	resp.Infos = infos

	return
}
