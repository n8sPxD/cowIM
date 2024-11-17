package logic

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/internal/business/auth/api/internal/svc"
	"github.com/n8sPxD/cowIM/internal/business/auth/api/internal/types"
	"github.com/n8sPxD/cowIM/internal/common/db/myMysql/models"
	"github.com/n8sPxD/cowIM/pkg/encrypt"
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
		logx.Error("[Register] Encrypt error:", err)
		return nil, errors.New("注册失败！好像是服务器发生了异常")
	}

	// 数据入库
	user := models.User{
		Username: req.Username,
		Password: password,
		Avatar:   req.Username[:0],
	}
	if err := l.svcCtx.MySQL.InsertUser(l.ctx, &user); err != nil {
		logx.Error("[Register] Insert user to DB failed, error: ", err)
		return nil, errors.New("注册失败！好像是服务器发生了异常")
	}
	return &types.RegisterResponse{ID: uint32(user.ID)}, nil
}
