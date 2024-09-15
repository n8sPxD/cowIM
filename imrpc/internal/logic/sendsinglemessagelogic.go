package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/common/db/models"
	"github.com/n8sPxD/cowIM/imrpc/imrpc"
	"github.com/n8sPxD/cowIM/imrpc/internal/svc"
	"google.golang.org/protobuf/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSingleMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSingleMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSingleMessageLogic {
	return &SendSingleMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendSingleMessageLogic) SendSingleMessage(in *imrpc.SendMessageRequest) (*imrpc.SendMessageResponse, error) {
	msg := models.SingleMessage{
		// TODO: 分配唯一递增消息id
		SendUser: in.SendUser,
		RecvUser: in.RecvUser,
		Type:     int8(in.Type),
		Content:  in.Content,
	}
	// 先入库
	if err := l.svcCtx.MySQL.InsertSingleMessage(l.ctx, &msg); err != nil {
		logx.Error("InsertSingleMessage error, err: ", err)
		return nil, err
	}
	// 再发配
	b, err := proto.Marshal(in)
	if err != nil {
		logx.Error("InsertMessage marshal error, err: ", err)
		return nil, err
	}

	// 发送到消息队列
	if err := l.svcCtx.KqPusherClient.Push(l.ctx, string(b)); err != nil {
		logx.Error("InsertMessage push message to Kafka error, err: ", err)
		return nil, err
	}

	return &imrpc.SendMessageResponse{}, nil
}
