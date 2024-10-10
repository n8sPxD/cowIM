package handler

import (
	"net/http"

	"github.com/n8sPxD/cowIM/common/response"
	"github.com/n8sPxD/cowIM/microservices/wsget/internal/logic"
	"github.com/n8sPxD/cowIM/microservices/wsget/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/wsget/internal/types"
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