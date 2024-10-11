package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/common/response"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/logic"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func groupInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupInfoRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpFail(r, w, err)
			return
		}

		l := logic.NewGroupInfoLogic(r.Context(), svcCtx)
		resp, err := l.GroupInfo(&req)
		if err != nil {
			response.HttpFail(r, w, err)
		} else {
			response.HttpSuccess(r, w, resp)
		}
	}
}
