package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/common/response"
	"github.com/n8sPxD/cowIM/microservices/auth/api/internal/logic"
	"github.com/n8sPxD/cowIM/microservices/auth/api/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/auth/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpFail(r, w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			response.HttpFail(r, w, err)
		} else {
			response.HttpSuccess(r, w, resp)
		}
	}
}
