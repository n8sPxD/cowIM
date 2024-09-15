package test

import (
	"fmt"
	"net"
	"testing"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/libnet"
)

func TestTcpConnChat(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	defer conn.Close()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	session := libnet.NewSession(nil, conn, 10)
	msg := libnet.Message{
		Header: libnet.Header{
			Command: constant.SINGLE_CHAT_REQ,
		},
		Body: []byte("哈哈"),
	}

	// 发送消息
	if err := session.Send(msg); err != nil {
		t.Error(err)
		t.Fail()
	}

	// 接受消息
	resp, err := session.Receive()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	fmt.Println(string(resp.Body))
}
