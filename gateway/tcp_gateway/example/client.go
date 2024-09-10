package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/n8sPxD/cowIM/common/protocol"
)

// TCP客户端结构体
type ClientConnect struct {
	conn *protocol.Connect
}

// NewClientConnect 创建客户端连接
func NewClientConnect(address string) (*ClientConnect, error) {
	conn, err := net.Dial("tcp", address) // 连接网关服务器
	if err != nil {
		return nil, err
	}
	return &ClientConnect{conn: &protocol.Connect{Conn: conn}}, nil
}

func main() {
	// 创建TCP客户端并连接到网关
	client, err := NewClientConnect("127.0.0.1:9000")
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer func() {
		_ = client.conn.Conn.Close()
	}()

	// 构建测试消息
	testMessage := protocol.Message{
		Header: protocol.Header{
			Version:      1,     // 假设版本号为1
			StatusCode:   200,   // 状态码200表示正常
			MessageType:  1,     // 假设消息类型为1
			Command:      100,   // 假设命令为100
			ClientSeqNum: 12345, // 客户端序列号
		},
		Body: []byte("Hello from client!"), // 消息体
	}

	// 使用已经封装的Send方法发送消息
	err = client.conn.Send(testMessage)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	fmt.Println("Message sent successfully!")

	// 使用封装的Receive方法接收服务器响应
	responseMsg, err := client.conn.Receive()
	if err != nil {
		log.Fatalf("Failed to receive message: %v", err)
	}

	// 打印接收到的消息
	fmt.Printf("Received message from server: %+v\n", responseMsg)
	fmt.Printf("Body: %s", responseMsg.Body)

	// 等待1秒后结束，模拟处理
	time.Sleep(1 * time.Second)
}
