// server/tcp_server.go
// 建立TCP服务器

package server

import (
	"context"
	"io"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/libnet"
	"github.com/n8sPxD/cowIM/common/socket"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/svc"
	"github.com/n8sPxD/cowIM/imrpc/imrpc"
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
		go s.sessionLoop(session)
	}
}

func (s *TcpServer) sessionLoop(ses *libnet.Session) {
	// 先接受Token进行认证
	//message, err := ses.Receive()
	//if err != nil {
	//	logx.Errorf("Receive token error, err: %v\n", err)
	//	ses.Close()
	//	return
	//}
	//token := string(message.Body)
	//name, err := jwt.ParseToken(token, s.svcCtx.Config.Auth.AccessSecret)
	//if err != nil {
	//	logx.Errorf("Parse JWT error, err: %v\n", err)
	//	ses.Close()
	//	return
	//}
	//// TODO: 同用户同时只能登陆一个客户端
	//ses.SetUser(name.PayLoad.Username)

	ses.SetUser("test")
	s.Server.Manager.Add("test", ses)

	// TODO: 持续接收消息，解析libnet.Message，根据Message.Cmd匹配不同的Rpc接口进行调用
	for {
		message, err := ses.Receive()
		if err != nil {
			if err == io.EOF {
				logx.Infof("User %s disconnect", ses.User())
				ses.Close()
				return
			}
			logx.Errorf("Receive message error, err: %v\n", err)
			ses.Close()
			return
		}
		if message.Command == constant.SINGLE_CHAT_REQ {
			// 测试，先给自己发消息
			req := imrpc.SendMessageRequest{
				SendUser:  ses.User(),
				RecvUser:  ses.User(),
				RecvGroup: constant.NONE_GROUP,
				Type:      constant.SINGLE_CHAT_REQ,
				Content:   string(message.Body),
			}

			if _, err := s.svcCtx.IMRpc.SendSingleMessage(context.Background(), &req); err != nil {
				logx.Error("Send RPC request error, err: ", err)
				ses.Close()
				return
			}
		}
	}
}
