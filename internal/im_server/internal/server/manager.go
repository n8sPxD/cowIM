package server

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

type Session struct {
	ID       UserID
	Username string
	Conn     *websocket.Conn
	mutex    sync.Mutex
}

type UserID uint32

type IConnectionManager interface {
	Add(*Session)
	Remove(uint32)
	RemoveWithCode(uint32, int, string)
	Get(uint32) (*Session, bool)
	SendMessage(uint32, []byte) error
	ReadMessage(uint32) ([]byte, error)
	GetAckHandler() IAckHandler
}

// ConnectionManager WebSocket连接管理器
type ConnectionManager struct {
	// TODO: map换ConcurrentMap
	connections map[UserID]*Session
	AckHandler  IAckHandler
	mutex       sync.RWMutex
}

func NewConnectionManager() IConnectionManager {
	return &ConnectionManager{
		connections: make(map[UserID]*Session),
		AckHandler:  NewAckHandler(),
		mutex:       sync.RWMutex{},
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

func (cm *ConnectionManager) RemoveWithCode(userID uint32, code int, err string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if s, ok := cm.connections[UserID(userID)]; ok {
		s.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, err))
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

var ClientGoingAway = errors.New("user is offline")

// SendMessage 服务器发送或转发消息 msg 给指定 userID
func (cm *ConnectionManager) SendMessage(userID uint32, msg []byte) error {
	s, ok := cm.Get(userID)
	if !ok {
		// 用户不在线
		return ClientGoingAway
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Conn.WriteMessage(websocket.BinaryMessage, msg)
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

func (cm *ConnectionManager) GetAckHandler() IAckHandler {
	return cm.AckHandler
}
