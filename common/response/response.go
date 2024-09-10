// response.go
// HttpSuccess, HttpFail 一般用于API接口响应
// HttpResponse 用于网关等

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

// HttpSuccess 成功响应
func HttpSuccess(r *http.Request, w http.ResponseWriter, resp any) {
	response := Resp{
		Code:    0,
		Msg:     "中嘞",
		Content: resp,
	}
	httpx.WriteJsonCtx(r.Context(), w, http.StatusOK, &response)
}

// HttpFail 失败响应
func HttpFail(r *http.Request, w http.ResponseWriter, err error) {
	response := Resp{
		Code:    6,
		Msg:     "寄了",
		Content: err.Error(),
	}
	httpx.WriteJsonCtx(r.Context(), w, http.StatusOK, &response)
}

// HttpResponse 自定义响应
func HttpResponse(r *http.Request, w http.ResponseWriter, status int, resp *Resp) {
	httpx.WriteJsonCtx(r.Context(), w, status, resp)
}
