package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/microservices/user/internal/logic"
	"github.com/n8sPxD/cowIM/microservices/user/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/user/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func addFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddFriendRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAddFriendsLogic(r.Context(), svcCtx)
		resp, err := l.AddFriends(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
