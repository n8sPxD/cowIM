package logic

import (
	"context"
	"errors"
	"strconv"

	"github.com/n8sPxD/cowIM/common/db/myMysql/models"
	"github.com/n8sPxD/cowIM/microservices/info/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/info/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChatListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatListLogic {
	return &ChatListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatListLogic) ChatList(req *types.ChatListRequest) (*types.ChatListResponse, error) {
	chatList, err := l.svcCtx.Mongo.GetRecentChatList(l.ctx, req.UserID, 10)
	if err != nil {
		// 从数据库查最新消息列表失败了，啥也不返回
		logx.Error("[ChatList] Get chat list from mongoDB failed, error: ", err)
		return nil, errors.New("获取对话列表失败")
	}
	var retList []types.ChatListInfo
	for _, chat := range chatList {
		var (
			info types.ChatListInfo
			user *models.User
		)
		user, err = l.svcCtx.MySQL.GetUserBaseInfo(l.ctx, chat.SenderID)
		if err != nil {
			// 从数据库查用户信息失败了，返回默认值
			logx.Error("[ChatList] Get user base info from MySQL failed, error: ", err)
			info.Username = strconv.FormatInt(int64(chat.SenderID), 10) // 用户的用户名为数字id
			info.RecentMsg = chat.RecentMsg
			info.Avatar = ""
		} else {
			info.Username = user.Username
			info.RecentMsg = chat.RecentMsg
			info.Avatar = user.Avatar
		}
		retList = append(retList, info)
	}

	return &types.ChatListResponse{List: retList}, nil
}
