package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/n8sPxD/cowIM/internal/common/constant"
	"github.com/n8sPxD/cowIM/internal/common/message/front"
	"google.golang.org/protobuf/proto"
)

func TestLoginAndSendMsg(t *testing.T) {
	// 步骤 1：用户登录并获取 JWT 令牌
	jwtToken, err := login(1, "123456")
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	fmt.Println("获取到的 JWT 令牌:", jwtToken)

	// 步骤 2：建立 WebSocket 连接并进行通信
	err = connectWebSocket(jwtToken)
	if err != nil {
		log.Fatalf("WebSocket 连接失败: %v", err)
	}
}

// LoginRequest 定义登录请求的结构
type LoginRequest struct {
	ID       uint32 `json:"id"`
	Password string `json:"password"`
}

// LoginResponse 定义登录响应的结构，包含嵌套的 content 字段
type LoginResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Content struct {
		Token string `json:"token"`
	} `json:"content"`
}

// login 执行登录操作并返回 JWT 令牌
func login(id uint32, password string) (string, error) {
	loginURL := "http://localhost:8080/login"

	// 创建登录请求数据
	loginReq := LoginRequest{
		ID:       id,
		Password: password,
	}

	// 将请求数据编码为 JSON
	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		return "", fmt.Errorf("无法编码登录请求: %v", err)
	}

	// 发送 HTTP POST 请求
	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("登录失败，状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应体
	var loginResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return "", fmt.Errorf("无法解析登录响应: %v", err)
	}

	// 检查返回的 code 是否为 0（假设 0 表示成功）
	if loginResp.Code != 0 {
		return "", fmt.Errorf("登录失败，消息: %s", loginResp.Msg)
	}

	return loginResp.Content.Token, nil
}

// connectWebSocket 建立 WebSocket 连接并发送/接收消息
func connectWebSocket(jwtToken string) error {
	// 定义 WebSocket URL，注意使用 "ws" 或 "wss" 协议
	wsURL := "ws://localhost:8081/ws"

	// 设置 WebSocket 连接的 HTTP 头，包含 JWT 令牌
	header := http.Header{}
	header.Set("Authorization", "Bearer "+jwtToken)

	// 建立 WebSocket 连接
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		return fmt.Errorf("无法建立 WebSocket 连接: %v", err)
	}
	defer conn.Close()
	fmt.Println("成功建立 WebSocket 连接")

	// 示例：发送一条消息
	messageToSend := "哈哈，这是一条测试消息😄"
	sendMessage := front.Message{
		From:      1,
		To:        2,
		Content:   messageToSend,
		Type:      constant.SINGLE_CHAT,
		MsgType:   constant.MSG_COMMON_MSG,
		Extend:    nil,
		Timestamp: time.Now().Unix(),
	}
	realMsg, err := proto.Marshal(&sendMessage)
	err = conn.WriteMessage(websocket.BinaryMessage, realMsg)
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}
	fmt.Printf("发送消息: %s\n", messageToSend)

	// 读取并打印服务器的响应
	_, message, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("读取消息失败: %v", err)
	}
	var msg front.Message
	err = proto.Unmarshal(message, &msg)
	if err != nil {
		return fmt.Errorf("读取消息失败: %v", err)
	}
	fmt.Printf("收到服务器响应: %s\n", msg.Content)

	return nil
}
