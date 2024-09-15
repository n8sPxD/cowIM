// socket/socket.go
// 建立一个服务器

package socket

import (
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/n8sPxD/cowIM/common/libnet"
)

type Server struct {
	Manager      *libnet.Manager // 连接管理器
	Listener     net.Listener    // 监听器
	Protocol     libnet.Protocol // 使用的通信协议
	SendChanSize int             // 消息队列缓冲区大小
}

func NewServer(l net.Listener, p libnet.Protocol, chanSize int) *Server {
	return &Server{
		Manager:      libnet.NewIMManager(),
		Listener:     l,
		Protocol:     p,
		SendChanSize: chanSize,
	}
}

// NewTcpServe 建立新的TCP服务器
func NewTcpServe(address string, protocol libnet.Protocol, sendChanSize int) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewServer(listener, protocol, sendChanSize), nil
}

func (ts *Server) Close() {
	ts.Listener.Close()
	ts.Manager.CloseAll()
}

func (ts *Server) Accept() (*libnet.Session, error) {
	var tempDelay time.Duration
	for {
		conn, err := ts.Listener.Accept()
		// 超时重传
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if deadline := 1 * time.Second; tempDelay > deadline {
					tempDelay = deadline
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				// 连接已关闭
				return nil, io.EOF
			}
			return nil, err
		}

		return libnet.NewSession(ts.Manager, conn, ts.SendChanSize), nil
	}
}
