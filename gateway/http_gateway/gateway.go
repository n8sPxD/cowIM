package main

import (
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/n8sPxD/cowIM/common/jwt"
	"github.com/n8sPxD/cowIM/common/response"
	"github.com/n8sPxD/cowIM/gateway/http_gateway/config"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

// Gateway 结构体，封装了 HTTP 网关的相关逻辑
type Gateway struct {
	server  *http.Server
	proxies []config.Proxy
}

// NewGateway 创建并初始化 HTTP 短连接网关
func NewGateway(proxies []config.Proxy) *Gateway {
	return &Gateway{proxies: proxies}
}

// Start 启动网关服务
func (g *Gateway) Start() error {
	mux := http.NewServeMux()

	// 设置反向代理路由
	for _, proxy := range g.proxies {
		mux.HandleFunc(proxy.Route, g.reverseProxyHandler(proxy.Target))
		logx.Infof("Route: %s -> Target: %s", proxy.Route, proxy.Target)
	}

	// 初始化 HTTP 服务器
	g.server = &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logx.Info("HTTP Gateway is running on :8080")
	return g.server.ListenAndServe()
}

func jwtParse(r *http.Request, w http.ResponseWriter) bool {
	// 提取 JWT Token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		response.HttpResponse(r, w, http.StatusUnauthorized, &response.Resp{
			Code:    6,
			Msg:     "寄了",
			Content: "请携带Token",
		})
		logx.Error("Authorization header is missing")
		return false
	}

	// 去掉 "Bearer " 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 解析JWT Token
	claims, err := jwt.ParseToken(tokenString, c.Auth.AccessSecret)
	if err != nil {
		response.HttpResponse(r, w, http.StatusUnauthorized, &response.Resp{
			Code:    6,
			Msg:     "寄了",
			Content: "Token解析失败",
		})
		logx.Errorf("Invalid token: %v", err)
		return false
	}

	// 鉴权成功，继续转发请求
	logx.Infof("Authenticated user: %s", claims.Username)
	return true
}

// reverseProxyHandler 返回一个 HTTP 处理器用于反向代理请求
func (g *Gateway) reverseProxyHandler(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 打印用户的请求地址
		logx.Infof("Received request: %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)

		// 排除白名单内的地址
		if _, ok := whiteList[r.URL.String()]; !ok {
			// 鉴权
			if !jwtParse(r, w) {
				return
			}
		}

		// 创建反向代理并转发请求
		backendURL, err := url.Parse(target)
		if err != nil {
			response.HttpResponse(r, w, http.StatusInternalServerError, &response.Resp{
				Code:    6,
				Msg:     "寄了",
				Content: "解析后端地址失败",
			})
			logx.Error("Failed to parse backend URL:", err)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(backendURL)
		proxy.ServeHTTP(w, r)
	}
}

// Shutdown 关闭网关服务
func (g *Gateway) Shutdown() error {
	logx.Info("Shutting down HTTP Gateway...")
	return g.server.Close()
}

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

var c config.Config
var whiteList = map[string]bool{}

func main() {
	flag.Parse()

	// 加载配置文件
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	// 初始化白名单
	for _, path := range c.WhiteList {
		whiteList[path] = true
	}

	// 初始化网关
	gateway := NewGateway(c.Proxies)

	// 处理退出信号，平滑关闭
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		_ = gateway.Shutdown()
		os.Exit(0)
	}()

	// 启动网关服务
	err := gateway.Start()
	if err != nil {
		logx.Error("Failed to start gateway: ", err)
	}
}
