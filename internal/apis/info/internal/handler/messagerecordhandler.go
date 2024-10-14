package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/internal/apis/info/internal/logic"
	"github.com/n8sPxD/cowIM/internal/apis/info/internal/svc"
	"github.com/n8sPxD/cowIM/internal/apis/info/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func messageRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessageRecordRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewMessageRecordLogic(r.Context(), svcCtx)
		resp, err := l.MessageRecord(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
