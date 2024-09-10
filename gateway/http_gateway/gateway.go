package main

import (
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/n8sPxD/cowIM/gateway/http_gateway/config"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

// 反向代理处理程序
func reverseProxy(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析目标 URL
		backendURL, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
			logx.Error("Invalid backend URL:", err)
			return
		}

		// 创建反向代理
		proxy := httputil.NewSingleHostReverseProxy(backendURL)

		// 转发请求
		proxy.ServeHTTP(w, r)
	}
}

var configFile = flag.String("f", "etc/proxy.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.Log)

	for _, info := range c.Proxies {
		http.HandleFunc(info.Route, reverseProxy(info.Target))
	}

	logx.Info("Gateway is running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logx.Error("ListenAndServe: ", err)
	}
}
