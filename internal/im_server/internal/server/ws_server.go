// server/ws_server.go
// 负责处理长连接的建立、保持以及消息的转发

package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/n8sPxD/cowIM/internal/common/constant"
	"github.com/n8sPxD/cowIM/internal/common/message/front"
	"github.com/n8sPxD/cowIM/internal/im_server/svc"
	"github.com/n8sPxD/cowIM/pkg/servicehub"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type Server struct {
	ctx         context.Context         // 上下文
	svcCtx      *svc.ServiceContext     // 依赖服务
	Manager     IConnectionManager      // 连接管理器
	upgrader    *websocket.Upgrader     // Websocket协议升级器
	messages    chan string             // 本地消息队列，作用是消息聚合
	close       chan struct{}           // 关闭信号
	registerHub *servicehub.RegisterHub // 注册中心
}

func MustNewServer(ctx context.Context, svcCtx *svc.ServiceContext) *Server {
	return &Server{
		ctx:     ctx,
		svcCtx:  svcCtx,
		Manager: NewConnectionManager(),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		messages:    make(chan string, 100000),
		close:       make(chan struct{}),
		registerHub: servicehub.NewRegisterHub(svcCtx.Config.Etcd.Hosts, 3),
	}
}

func (s *Server) Start() {
	// 创建路由
	r := mux.NewRouter()
	r.HandleFunc("/ws", s.handleWebSocket)

	// 注册服务
	s.register()

	// 启动HTTP服务器
	addr := fmt.Sprintf("%s:%d", s.svcCtx.Config.Host, s.svcCtx.Config.Port)
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
	// TODO: 责任链模式重构handleWebsocket
	if online := s.checkOnline(id); online {
		return
	}

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
	logx.Infof("[handleWebsocket] User %s on %s connected", name, conn.RemoteAddr())

	// 维护用户登陆路由
	if _, err := s.svcCtx.Redis.UpdateUserRouterStatus(s.ctx, id, s.svcCtx.Config.WorkID); err != nil {
		logx.Error("[updateRouterStatus] Update router status to redis failed, error: ", err)
	}

	// 读消息
	go func() {
		err = s.readMessageFromFrontend(id)
		if err != nil {
			logx.Debug("[handleWebsocket] Read message failed, error: ", err)
			s.close <- struct{}{}
		}
	}()
	// 发送到消息队列处理
	go func() {
		err = s.sendMessageToBackend()
		if err != nil {
			logx.Debug("[handleWebsocket] Send message failed, error: ", err)
			s.close <- struct{}{}
		}
	}()

	for {
		<-s.close
		s.svcCtx.Redis.RemoveUserRouterStatus(
			s.ctx,
			id,
		) // TODO: 这里出错要么是没读到在线状态，不影响，要么就是redis挂了，再考虑后续怎么处理
		return
	}
}

func (s *Server) checkOnline(id uint32) bool {
	// 处理重复登陆, 把已经登陆的客户端踢下线

	// 先从本地找，再从redis找
	if _, online := s.Manager.Get(id); online {
		// 在当前服务器在线, 直接踢
		s.Manager.RemoveWithCode(id, constant.DUP_CLIENT_CODE, constant.DUP_CLIENT_ERR)
		return true
	}

	if _, err := s.svcCtx.Redis.GetUserRouterStatus(s.ctx, id); err == nil {
		// 没出错，说明找到了用户在线，但是不在当前服务器中
		// 消息塞队列，给隔壁处理
		if msg, err := proto.Marshal(&front.Message{
			From:    constant.USER_SYSTEM,
			To:      id,
			Type:    constant.SYSTEM_INFO,
			MsgType: constant.MSG_DUP_CLIENT,
		}); err != nil {
			logx.Error("[checkOnline] Marshal message to protobuf failed, error: ", err)
			return false
		} else {
			s.messages <- string(msg)
			return true
		}
	}
	return false
}
