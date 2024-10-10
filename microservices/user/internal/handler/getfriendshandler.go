package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/common/response"
	"github.com/n8sPxD/cowIM/microservices/user/internal/logic"
	"github.com/n8sPxD/cowIM/microservices/user/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/user/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFriendsRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpFail(r, w, err)
			return
		}

		l := logic.NewGetFriendsLogic(r.Context(), svcCtx)
		resp, err := l.GetFriends(&req)
		if err != nil {
			response.HttpFail(r, w, err)
		} else {
			response.HttpSuccess(r, w, resp)
		}
	}
}
