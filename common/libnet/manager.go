// libnet/manager.go
// 管理连接进入的TCP TcpSessions，参考了 https://github.com/zhoushuguang/zeroim/tree/main/common/libnet (直接说照着抄吧)

package libnet

import (
	"errors"
	"sync"
)

var TcpSessionNotFound = errors.New("user not login")

type Manager struct {
	Sessions map[string]*Session // key: 用户token
	sync.RWMutex

	disposeFlag bool
	disposeOnce sync.Once
	disposeWait sync.WaitGroup
}

func NewIMManager() *Manager {
	return &Manager{
		Sessions: make(map[string]*Session),
	}
}

func (m *Manager) Add(username string, session *Session) {
	m.Lock()
	defer m.Unlock()
	m.Sessions[username] = session
}

func (m *Manager) Remove(username string) {
	m.Lock()
	defer m.Unlock()
	if !m.Sessions[username].IsClosed() {
		m.Sessions[username].Close()
	}
	delete(m.Sessions, username)

}

func (m *Manager) Get(username string) (*Session, error) {
	m.RLock()
	defer m.RUnlock()
	s, exists := m.Sessions[username]
	if !exists {
		return nil, TcpSessionNotFound
	}
	return s, nil
}

func (m *Manager) CloseAll() {
	m.disposeOnce.Do(func() {
		m.disposeFlag = true
		for n := range m.Sessions {
			session := m.Sessions[n]
			m.Lock()
			_ = session.Close()
			m.Unlock()
		}
		m.disposeWait.Wait()
	})
}
