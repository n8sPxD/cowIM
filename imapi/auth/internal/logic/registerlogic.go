package logic

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/n8sPxD/cowIM/common/db/models"
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
	if err := l.svcCtx.MySQL.InsertUser(l.ctx, &models.User{
		Username: req.Nickname,
		Password: req.Password,
	}); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, errors.New("该用户名已被注册")
		}
		return nil, err
	}

	return &types.RegisterResponse{}, nil
}
