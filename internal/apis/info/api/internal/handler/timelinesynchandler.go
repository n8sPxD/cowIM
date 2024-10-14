package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/internal/apis/info/api/internal/logic"
	"github.com/n8sPxD/cowIM/internal/apis/info/api/internal/svc"
	"github.com/n8sPxD/cowIM/internal/apis/info/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func timelineSyncHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TimelineSyncRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewTimelineSyncLogic(r.Context(), svcCtx)
		resp, err := l.TimelineSync(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
