package internal

import (
	"sync"

	"github.com/gorilla/websocket"
	__front "github.com/n8sPxD/cowIM/common/message/.front"
	"google.golang.org/protobuf/proto"
)

type Session struct {
	ID       UserID
	Username string
	Conn     *websocket.Conn
	mutex    sync.Mutex
}

type UserID uint32

// ConnectionManager WebSocket连接管理器
type ConnectionManager struct {
	// TODO: map换ConcurrentMap
	connections map[UserID]*Session
	mutex       sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[UserID]*Session),
	}
}

func (cm *ConnectionManager) Add(s *Session) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connections[s.ID] = s
}

func (cm *ConnectionManager) Remove(userID uint32) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if s, ok := cm.connections[UserID(userID)]; ok {
		s.Conn.Close()
		delete(cm.connections, UserID(userID))
	}
}

func (cm *ConnectionManager) Get(userID uint32) (*Session, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	s, ok := cm.connections[UserID(userID)]
	return s, ok
}

// SendMessage 服务器发送或转发消息 msg 给指定 userID
func (cm *ConnectionManager) SendMessage(userID uint32, msg *__front.Message) error {
	cm.mutex.RLock()
	s, ok := cm.connections[UserID(userID)]
	cm.mutex.RUnlock()
	if !ok {
		// 用户不在线
		return nil
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.Conn.WriteMessage(websocket.BinaryMessage, data)
}

// ReadMessage 服务器从 userID 接受消息
func (cm *ConnectionManager) ReadMessage(userID uint32) ([]byte, error) {
	msgType, msg, err := cm.connections[UserID(userID)].Conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if msgType != websocket.BinaryMessage {
		return nil, err
	}
	return msg, nil
}
