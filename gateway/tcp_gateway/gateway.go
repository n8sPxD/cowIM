package main

import (
	"flag"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/n8sPxD/cowIM/common/protocol"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/config"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

// Gateway 长连接网关结构
type Gateway struct {
	listener net.Listener
}

// NewGateway 创建一个新的TCP服务器
func NewGateway(address string) (*Gateway, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Gateway{listener: listener}, nil
}

// Start 启动网关，开始接受客户端连接
func (g *Gateway) Start() {
	logx.Info("Gateway started, waiting for connections...")
	for {
		conn, err := g.listener.Accept()
		if err != nil {
			logx.Errorf("Failed to accept connection: %v", err)
			continue
		}
		go g.handleConnection(conn) // 每个连接交由一个goroutine处理
	}
}

// handleConnection 处理单个TCP连接
func (g *Gateway) handleConnection(conn net.Conn) {
	logx.Infof("Accepted connection from %s\n", conn.RemoteAddr().String())
	client := &protocol.Connect{Conn: conn}

	defer func() {
		_ = conn.Close()
	}()

	// 接收消息
	for {
		msg, err := client.Receive() // 读取并解码消息
		if err != nil {
			if err == io.EOF {
				logx.Infof("User disconnected from gateway")
			} else {
				logx.Errorf("Error receiving message: %v", err)
			}
			return
		}

		// 处理消息中的 Command 字段，调用RPC接口
		logx.Debugf("Received Message: %+v", msg)
		// 这里可以根据msg.Command调用RPC接口，然后处理返回的结果

		// 返回一个简单的响应消息
		// TODO: 根据用户请求以及长连接后端服务响应进行改造
		response := protocol.Message{
			Header: protocol.Header{
				Version:      msg.Version,
				StatusCode:   200, // 假设返回成功状态码
				MessageType:  msg.MessageType,
				Command:      msg.Command,
				ClientSeqNum: msg.ClientSeqNum,
			},
			Body: []byte("Hello from n8spxd's gateway!!!!"),
		}

		if err := client.Send(response); err != nil {
			logx.Errorf("Error sending response: %v", err)
		}
	}
}

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.Log)

	gateway, err := NewGateway(":9000") // 监听9000端口
	if err != nil {
		logx.Errorf("Failed to start gateway: %v", err)
	}

	// 处理退出信号，平滑关闭
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		logx.Info("Shutting down gateway...")
		_ = gateway.listener.Close()
		os.Exit(0)
	}()

	// 启动网关
	gateway.Start()
}
