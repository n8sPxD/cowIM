// server/websocket.go
// 负责处理长连接的建立、保持以及消息的转发

package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/n8sPxD/cowIM/common/jwt"
	"github.com/n8sPxD/cowIM/im-server/internal"
	"github.com/n8sPxD/cowIM/im-server/internal/config"
	"github.com/n8sPxD/cowIM/im-server/svc"
	"github.com/n8sPxD/cowIM/microservices/auth/rpc/types/authRpc"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	config   config.Config               // Server的设置
	Manager  *internal.ConnectionManager // 连接管理器
	upgrader *websocket.Upgrader         // Websocket协议升级器
	ctx      context.Context             // 上下文
	svcCtx   *svc.ServiceContext         // 依赖服务
}

func MustNewServer(c config.Config, ctx context.Context, svcCtx *svc.ServiceContext) *Server {
	return &Server{
		config:  c,
		Manager: internal.NewConnectionManager(),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (s *Server) Start() {
	// 创建路由
	r := mux.NewRouter()
	r.HandleFunc("/ws", s.handleWebSocket)

	// 启动HTTP服务器
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	logx.Info("WebSocket server starting on ", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		logx.Error("ListenAndServe: ", err)
	}
}

func (s *Server) Close() {
}

// handleWebSocket 处理WebSocket连接
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 先鉴权，再进行消息通讯
	// 从Authorization头部获取JWT令牌
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	// JWT 格式： Bearer <token>
	var tokenString string
	fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
	if tokenString == "" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	// 验证 JWT
	claims, err := jwt.ParseToken(tokenString, s.config.Auth.AccessSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// -------------- JWT 验证 ------------------

	type User struct {
		ID       uint32
		Username string
	}
	user := &User{
		ID:       claims.PayLoad.ID,
		Username: claims.PayLoad.Username,
	}

	// TODO: 完善重复登陆逻辑
	// 处理重复登陆, 把已经登陆的客户端踢下线
	/*
		if session, online := s.Manager.Get(user.ID); online {
			msg := message.Message{
				ID:        idgen.NextId(),
				From:      constant.USER_SYSTEM,
				To:        user.ID, // 发给原客户端
				Content:   "你已经在另一个客户端上登陆",
				Type:      constant.MSG_SYSTEM_MSG,
				Extend:    nil,
				Timestamp: uint64(time.Now().Unix()),
			}
			protoMsg, err := proto.Marshal(&msg)
			if err != nil {
				logx.Error("[handleWebsocket] Marshal message failed, error: ", err)
				return
			}
			// 此处conn为已经建立好的连接
			err = session.Conn.WriteMessage(websocket.TextMessage, protoMsg)
			if err != nil {
				logx.Error("[handleWebsocket] Write message failed, error: ", err)
				return
			}
		}
	*/

	// 升级 Websocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Infof("[handleWebsocket] Upgrade failed: %v", err)
		return
	}

	// 添加连接到管理器
	currentSession := internal.Session{
		ID:       internal.UserID(user.ID),
		Username: user.Username,
		Conn:     conn,
	}
	s.Manager.Add(&currentSession)
	defer s.Manager.Remove(user.ID)
	logx.Infof("[handleWebsocket] User %s connected", user.Username)

	// 维护用户登陆路由
	go func() {
		_, err = s.svcCtx.AuthRpc.UserConnStatus(s.ctx, &authRpc.ConnRequest{
			WorkId: uint32(s.config.WorkID),
			UserId: user.ID,
		})
		if err != nil {
			logx.Error("[handleWebsocket] RPC UserConnStatus failed, error: ", err)
		}
	}()

	// 处理消息
	for {
		// 读消息
		// 此处的消息为protobuf序列化后的消息
		msg, err := s.Manager.ReadMessage(user.ID)
		if err != nil {
			// 用户断线
			if websocket.IsCloseError(
				err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				logx.Infof("User %d disconnected", user.ID)

				// 用户非正常无法读取信息
			} else {
				logx.Errorf(
					"[handleWebsocket] Read message from user %v failed, error: %v",
					user.ID,
					err,
				)
			}
			logx.Debugf("[handleWebsocket] Removing user %d router status...", user.ID)
			s.svcCtx.Redis.RemoveUserRouterStatus(s.ctx, user.ID)
			return
		}

		// 发送到消息队列处理
		mqMsg := kafka.Message{
			Value: msg,
		}
		logx.Debug("[handleWebsocket] Pushing message to MQ...")
		if err := s.svcCtx.MsgForwarder.WriteMessages(s.ctx, mqMsg); err != nil {
			logx.Error("[handleWebsocket] Push message to MQ failed, error: ", err)
			return
		}
		logx.Debug("[handleWebsocket] Pushing over")
	}
}
