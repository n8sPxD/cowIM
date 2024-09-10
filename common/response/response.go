package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Resp struct {
	Code    uint8  `json:"code"`
	Msg     string `json:"msg"`
	Content any    `json:"content"`
}

func HttpSuccess(r *http.Request, w http.ResponseWriter, resp any) {
	response := Resp{
		Code:    0,
		Msg:     "中嘞",
		Content: resp,
	}
	httpx.WriteJsonCtx(r.Context(), w, http.StatusOK, &response)
}

func HttpFail(r *http.Request, w http.ResponseWriter, code uint8, err error) {
	response := Resp{
		Code:    code,
		Msg:     err.Error(),
		Content: nil,
	}
	httpx.WriteJsonCtx(r.Context(), w, http.StatusOK, &response)
}
