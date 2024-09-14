package main

import (
	"fmt"
	"io"
	"net"

	"github.com/n8sPxD/cowIM/common/libnet"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	protocol := libnet.NewIMProtocol()
	parser := protocol.NewParser()
	msg := parser.Encode(libnet.Message{
		Header: libnet.Header{
			Version:      50,
			StatusCode:   200,
			MessageType:  233,
			Command:      233,
			ClientSeqNum: 1000,
		},
		Body: []byte("我嘞个豆"),
	})
	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println(err)
	}

	b := make([]byte, 26)
	_, err = io.ReadFull(conn, b)
	resp, err := parser.Decode(b, 26)
	if err != nil {
		fmt.Println("错误啦, err: ", err)
	}
	fmt.Println(resp.Body)

	//reader := bufio.NewReader(conn)
	//// 读取直到换行符
	//response, err := reader.ReadString('\n')
	//if err != nil {
	//	fmt.Println("读取消息错误:", err)
	//}
	//fmt.Println("收到服务器消息:", response)
}
