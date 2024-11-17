package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/internal/business/info/api/internal/svc"
	"github.com/n8sPxD/cowIM/internal/business/info/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type MessageRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageRecordLogic {
	return &MessageRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageRecordLogic) MessageRecord(req *types.MessageRecordRequest) (resp *types.MessageRecordResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
