// server/tcp_server.go
// 建立TCP服务器

package server

import (
	"github.com/n8sPxD/cowIM/common/libnet"
	"github.com/n8sPxD/cowIM/common/socket"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type TcpServer struct {
	Server *socket.Server      // 服务端
	svcCtx *svc.ServiceContext // 服务依赖
}

func MustNewTcpServer(svcCtx *svc.ServiceContext, addr string) *TcpServer {
	protocol := libnet.NewIMProtocol()
	serv, err := socket.NewTcpServe(addr, protocol, 100)
	if err != nil {
		panic(err)
	}
	return &TcpServer{
		Server: serv,
		svcCtx: svcCtx,
	}
}

func (s *TcpServer) HandleRequest() {
	for {
		session, err := s.Server.Accept()
		if err != nil {
			panic(err)
		}
		session.User = "test"
		s.Server.Manager.Add("test", session)
		go s.sessionLoop(session)
	}
}

func (s *TcpServer) sessionLoop(ses *libnet.Session) {
	for {
		message, err := ses.Receive()
		if err != nil {
			logx.Errorf("Receive message error, err: %v\n", err)
			_ = ses.Close()
			return
		}
		logx.Infof("Receive message: %v, body: %s, packLen: %v\n", message.Header, string(message.Body))
		logx.Infof("sending message: %v ...", message)
		err = ses.Send(*message)
		if err != nil {
			logx.Errorf("Send message error, err: %v\n", err)
			_ = ses.Close()
			return
		}
	}
}
