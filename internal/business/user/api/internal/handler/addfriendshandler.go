package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/internal/business/user/api/internal/logic"
	"github.com/n8sPxD/cowIM/internal/business/user/api/internal/svc"
	"github.com/n8sPxD/cowIM/internal/business/user/api/internal/types"
	"github.com/n8sPxD/cowIM/pkg/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func addFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddFriendRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpFail(r, w, err)
			return
		}

		l := logic.NewAddFriendsLogic(r.Context(), svcCtx)
		resp, err := l.AddFriends(&req)
		if err != nil {
			response.HttpFail(r, w, err)
		} else {
			response.HttpSuccess(r, w, resp)
		}
	}
}
