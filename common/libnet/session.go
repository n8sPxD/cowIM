// libnet/session.go
// 管理单个TCP连接会话， 参考了 https://github.com/zhoushuguang/zeroim/tree/main/common/libnet (直接说照着抄吧)

package libnet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/logx"
)

// Session TCP会话
type Session struct {
	User       string       // 哪位用户连接的
	Conn       net.Conn     // 会话
	parser     Parser       // 协议解析器
	manager    *Manager     // 该会话属于哪个manager管理
	sendChan   chan Message // 消息缓存队列
	closeFlag  int32
	closeChan  chan int
	closeMutex sync.Mutex
}

func TcpSession(manager *Manager, Conn net.Conn, chanSize int) *Session {
	s := &Session{
		Conn:      Conn,
		parser:    NewIMProtocol().NewParser(),
		manager:   manager,
		closeFlag: 0,
		closeChan: make(chan int),
	}
	if chanSize > 0 {
		s.sendChan = make(chan Message, chanSize)
		go func() { s.sendLoop() }()
	}
	return s
}

func (ts *Session) sendLoop() {
	for {
		select {
		case msg := <-ts.sendChan:
			err := ts.send(msg)
			if err != nil {
				logx.Errorf("User %v send message error: %v", ts.User, err)
				return
			}
		case <-ts.closeChan:
			ts.Close()
			return
		}
	}
}

func (ts *Session) send(msg Message) error {
	buf := ts.parser.Encode(msg)
	n, err := ts.Conn.Write(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		logx.Errorf("n: %d, len(buf):%d\n", n, len(buf))
		return fmt.Errorf("发送失败！可能是服务器出了问题")
	}
	logx.Info("send over")
	return nil
}

var (
	SessionClosedErr  = errors.New("session closed")
	SessionBlockedErr = errors.New("session blocked")
)

func (ts *Session) Send(msg Message) error {
	if ts.IsClosed() {
		return SessionClosedErr
	}
	if ts.sendChan == nil {
		return ts.send(msg)
	}
	select {
	case ts.sendChan <- msg:
		return nil
	default:
		return SessionBlockedErr
	}
}

func (ts *Session) IsClosed() bool {
	return atomic.LoadInt32(&ts.closeFlag) == 1
}

func (ts *Session) readPackSize() (uint32, error) {
	return ts.readUint32BE()
}

func (ts *Session) readUint32BE() (uint32, error) {
	b := make([]byte, PACK_SIZE)
	_, err := io.ReadFull(ts.Conn, b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

func (ts *Session) readPacket(msgSize uint32) ([]byte, error) {
	b := make([]byte, msgSize)
	_, err := io.ReadFull(ts.Conn, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (ts *Session) Receive() (*Message, error) {
	packLen, err := ts.readPackSize()
	if err != nil {
		return nil, err
	}
	logx.Info("packLen: ", packLen)
	if packLen > MAX_PACK_SIZE {
		// TODO: 分包发送过长消息
		return nil, errors.New("发送消息过长")
	}
	buf, err := ts.readPacket(packLen)
	if err != nil {
		return nil, err
	}

	msg, err := ts.parser.Decode(buf, packLen)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ts *Session) Close() error {
	if atomic.CompareAndSwapInt32(&ts.closeFlag, 0, 1) {
		err := ts.Conn.Close()
		close(ts.closeChan)
		if ts.manager != nil {
			ts.manager.Remove(ts.User)
		}
		return err
	}

	return SessionClosedErr
}
