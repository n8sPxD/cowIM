package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/microservices/group/api/internal/logic"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func groupJoinedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupJoinedRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGroupJoinedLogic(r.Context(), svcCtx)
		resp, err := l.GroupJoined(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
