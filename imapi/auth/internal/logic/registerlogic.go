package logic

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/common/db/models"
	"github.com/n8sPxD/cowIM/common/encrypt"
	"github.com/n8sPxD/cowIM/imapi/auth/internal/svc"
	"github.com/n8sPxD/cowIM/imapi/auth/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (*types.RegisterResponse, error) {
	// Hash密码
	password, err := encrypt.HashPassword(req.Password)
	if err != nil {
		logx.Error("Encrypt error:", err)
		return nil, errors.New("注册失败！好像是服务器发生了异常")
	}

	// 数据入库，判断用户重名
	if err := l.svcCtx.MySQL.InsertUser(l.ctx, &models.User{
		Username: req.Username,
		Password: password,
	}); err != nil {
		return nil, err
	}
	return &types.RegisterResponse{}, nil
}
