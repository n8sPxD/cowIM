package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/internal/apis/wsget/api/internal/logic"
	"github.com/n8sPxD/cowIM/internal/apis/wsget/api/internal/svc"
	"github.com/n8sPxD/cowIM/internal/apis/wsget/api/internal/types"
	"github.com/n8sPxD/cowIM/pkg/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getAvailableWSServerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WebsocketServerGetRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpFail(r, w, err)
			return
		}

		l := logic.NewGetAvailableWSServerLogic(r.Context(), svcCtx)
		resp, err := l.GetAvailableWSServer(&req)
		if err != nil {
			response.HttpFail(r, w, err)
		} else {
			response.HttpSuccess(r, w, resp)
		}
	}
}
