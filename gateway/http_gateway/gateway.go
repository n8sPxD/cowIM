package main

import (
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// reverseProxyHandler 返回一个 HTTP 处理器用于反向代理请求
func (g *Gateway) reverseProxyHandler(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 打印用户的请求地址
		logx.Infof("Received request: %s %s from %s", r.Method, r.URL.String(), r.RemoteAddr)

		backendURL, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
			logx.Error("Failed to parse backend URL:", err)
			return
		}

		// 创建反向代理并转发请求
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

func main() {
	// 加载配置文件路径
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	// 初始化网关
	gateway := NewGateway(c.Proxies)

	// 处理退出信号，平滑关闭
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		logx.Info("Shutting down gateway...")
		_ = gateway.Shutdown()
		os.Exit(0)
	}()

	// 启动网关服务
	err := gateway.Start()
	if err != nil {
		logx.Error("Failed to start gateway: ", err)
	}
}
