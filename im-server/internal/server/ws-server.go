// server/ws-server.go
// 负责处理长连接的建立、保持以及消息的转发

package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/message/front"
	"github.com/n8sPxD/cowIM/im-server/internal/config"
	"github.com/n8sPxD/cowIM/im-server/svc"
	"github.com/n8sPxD/cowIM/microservices/auth/rpc/types/authRpc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	ctx      context.Context     // 上下文
	svcCtx   *svc.ServiceContext // 依赖服务
	config   config.Config       // Server的设置
	Manager  *ConnectionManager  // 连接管理器
	upgrader *websocket.Upgrader // Websocket协议升级器
	messages chan string         // 本地消息队列，作用是消息聚合
}

func MustNewServer(c config.Config, ctx context.Context, svcCtx *svc.ServiceContext) *Server {
	return &Server{
		ctx:     ctx,
		svcCtx:  svcCtx,
		config:  c,
		Manager: NewConnectionManager(),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		messages: make(chan string, 100000),
	}
}

func (s *Server) Start() {
	// 创建路由
	r := mux.NewRouter()
	r.HandleFunc("/ws", s.handleWebSocket)

	// 启动HTTP服务器
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	fmt.Println("WebSocket server starting on ", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		logx.Error("ListenAndServe: ", err)
	}
}

func (s *Server) Close() {
}

// handleWebSocket 处理WebSocket连接
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// JWT 鉴权
	var (
		id   uint32
		name string
	)
	if claims, ok := s.authenticate(w, r); ok {
		id = claims.PayLoad.ID
		name = claims.PayLoad.Username
	} else {
		return
	}

	// 处理重复在线
	s.checkOnline(id)

	// 升级 Websocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Error("[handleWebsocket] Upgrade to websocket failed, error: ", err)
		return
	}

	// 添加连接到管理器
	s.Manager.Add(&Session{
		ID:       UserID(id),
		Username: name,
		Conn:     conn,
	})
	defer s.Manager.Remove(id)
	logx.Infof("[handleWebsocket] User %s connected", name)

	// 维护用户登陆路由
	go s.updateRouterStatus(id)

	// 处理消息
	for {
		// 读消息
		err := s.readMessageFromFrontend(id)
		if err != nil {
			return
		}
		// 发送到消息队列处理
		err = s.sendMessageToBackend()
		if err != nil {
			return
		}
	}
}

func (s *Server) checkOnline(id uint32) {
	// 处理重复登陆, 把已经登陆的客户端踢下线
	// 给之前登陆的客户端的消息体
	var (
		msg []byte
		err error
	)
	if msg, err = proto.Marshal(&front.Message{
		From:    constant.SYSTEM,
		To:      id,
		Content: "您已在另一台客户端登陆！即将强制下线",
		Type:    constant.SYSTEM_INFO,
		MsgType: constant.MSG_SYSTEM_MSG,
		Extend:  nil,
		Time:    time.Now().Unix(),
	}); err != nil {
		logx.Error("[checkOnline] Marshal message to protobuf failed, error: ", err)
		// 发不了通知消息不影响后续把人家踢下线的流程，大不了之前的客户端干啥都干不了(连接都断了)，所以无需return
	}

	sendMsg := func() {
		// 判断msg的长度，防止上面marshal后没结果
		if len(msg) > 0 {
			s.messages <- string(msg)
		}
	}

	// 先从本地找，再从redis找
	_, online := s.Manager.Get(id)
	if online {
		// 在当前服务器在线, 直接发消息
		sendMsg()
		s.Manager.Remove(id)
	}

	_, err = s.svcCtx.Redis.GetUserRouterStatus(s.ctx, id)
	if err == nil {
		// 没出错，说明找到了用户在线
		sendMsg()
		s.Manager.Remove(id)
	}
}

func (s *Server) updateRouterStatus(id uint32) {
	_, err := s.svcCtx.AuthRpc.UserConnStatus(s.ctx, &authRpc.ConnRequest{
		WorkId: uint32(s.config.WorkID),
		UserId: id,
	})
	if err != nil {
		logx.Error("[updateRouterStatus] RPC UserConnStatus failed, error: ", err)
	}
}
