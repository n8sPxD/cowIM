package main

import (
	"errors"
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/n8sPxD/cowIM/internal/http_gateway/config"
	"github.com/n8sPxD/cowIM/pkg/response"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

// PayLoad 定义 JWT 中包含的用户信息
type PayLoad struct {
	ID       uint32 `json:"ID"`
	Username string `json:"username"`
}

// CustomClaims 定义自定义的 JWT Claims
type CustomClaims struct {
	PayLoad
	jwt.RegisteredClaims
}

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
		// 使用handleCORS包装每个处理器
		mux.HandleFunc(
			proxy.Route,
			handleCORS(http.HandlerFunc(g.reverseProxyHandler(proxy.Target))),
		)
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

// Data 存储解析后的用户信息
type Data struct {
	ID       uint32
	Username string
}

// jwtParse 解析 JWT Token 并返回用户数据
func jwtParse(r *http.Request, w http.ResponseWriter) (*Data, bool) {
	// 提取 JWT Token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		response.HttpResponse(r, w, http.StatusUnauthorized, &response.Resp{
			Code:    6,
			Msg:     "请携带Token",
			Content: "Authorization header is missing",
		})
		logx.Error("Authorization header is missing")
		return nil, false
	}

	// 去掉 "Bearer " 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 使用自定义的 ParseToken 函数解析 JWT Token
	claims, err := ParseToken(tokenString, c.Auth.AccessSecret)
	if err != nil {
		response.HttpResponse(r, w, http.StatusUnauthorized, &response.Resp{
			Code:    6,
			Msg:     "Token解析失败",
			Content: err.Error(),
		})
		logx.Errorf("Invalid token: %v", err)
		return nil, false
	}

	// 鉴权成功，返回用户数据
	logx.Infof(
		"Authenticated user ID: %d, name: %s",
		claims.PayLoad.ID,
		claims.PayLoad.Username,
	)
	return &Data{
		ID:       claims.PayLoad.ID,
		Username: claims.PayLoad.Username,
	}, true
}

// handleCORS 设置CORS响应头并处理OPTIONS预检请求
func handleCORS(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS响应头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, Authorization, UserID, Username")

		// 如果是OPTIONS请求，提前返回
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 否则继续处理下一个处理器
		next.ServeHTTP(w, r)
	}
}

// reverseProxyHandler 返回一个 HTTP 处理器用于反向代理请求
func (g *Gateway) reverseProxyHandler(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 打印用户的请求地址
		logx.Infof(
			"[reverseProxyHandler] Received request: %s %s from %s",
			r.Method,
			r.URL.String(),
			r.RemoteAddr,
		)

		// 排除白名单内的地址
		if _, ok := whiteList[r.URL.String()]; !ok {
			// 鉴权
			userData, authenticated := jwtParse(r, w)
			if !authenticated {
				return
			}

			// 在请求头中添加用户ID和用户名
			if userData != nil {
				r.Header.Set("UserID", strconv.FormatUint(uint64(userData.ID), 10))
				r.Header.Set("Username", userData.Username)
			}
		}

		// 创建反向代理并转发请求
		backendURL, err := url.Parse(target)
		if err != nil {
			response.HttpResponse(r, w, http.StatusInternalServerError, &response.Resp{
				Code:    6,
				Msg:     "解析后端地址失败",
				Content: err.Error(),
			})
			logx.Error("[reverseProxyHandler] Failed to parse backend URL:", err)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(backendURL)

		// 修改代理请求的请求头
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			// 头部已经在原始请求中设置，无需额外操作
		}

		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, e error) {
			response.HttpResponse(r, rw, http.StatusBadGateway, &response.Resp{
				Code:    6,
				Msg:     "代理请求失败",
				Content: e.Error(),
			})
			logx.Errorf("[reverseProxyHandler] Proxy error: %v", e)
		}

		proxy.ServeHTTP(w, r)
	}
}

// Shutdown 关闭网关服务
func (g *Gateway) Shutdown() error {
	logx.Info("[Gateway.Shutdown]Shutting down HTTP Gateway...")
	return g.server.Close()
}

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

var (
	c         config.Config
	whiteList = map[string]bool{}
)

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

// ParseToken 解析JWT token
func ParseToken(tokenString, accessSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(accessSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token非法")
}
