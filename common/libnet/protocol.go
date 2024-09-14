// protocol.go
// 自定义应用层协议，参考了 https://github.com/zhoushuguang/zeroim/tree/main/common/libnet (直接说照着抄吧)

package libnet

import (
	"encoding/binary"
	"fmt"
)

// 定义常量用于字段长度
const (
	MAX_BODY_SIZE = 1 << 12 // 单个包体最大长度

	// 字段长度
	PACK_SIZE           = 4                                                                                        // 总长度字段长度
	HEADER_SIZE         = 4                                                                                        // 消息头总长度字段长度
	VERSION_SIZE        = 1                                                                                        // 版本号字段长度
	STATUS_CODE_SIZE    = 1                                                                                        // 状态码字段长度
	MESSAGE_TYPE_SIZE   = 2                                                                                        // 消息类型字段长度
	COMMAND_SIZE        = 2                                                                                        // 命令字段长度
	CLIENT_SEQ_NUM_SIZE = 4                                                                                        // 客户端序列号字段长度
	TOTAL_HEADER_SIZE   = VERSION_SIZE + STATUS_CODE_SIZE + MESSAGE_TYPE_SIZE + COMMAND_SIZE + CLIENT_SEQ_NUM_SIZE // 头部长度
	MAX_PACK_SIZE       = MAX_BODY_SIZE + TOTAL_HEADER_SIZE + HEADER_SIZE + PACK_SIZE

	// 偏移量
	HEADER_OFFSET         = 0                                           // 头部长度字段偏移量
	VERSION_OFFSET        = HEADER_OFFSET + HEADER_SIZE                 // 版本号字段偏移量
	STATUS_OFFSET         = VERSION_OFFSET + VERSION_SIZE               // 状态码字段偏移量
	MESSAGE_TYPE_OFFSET   = STATUS_OFFSET + STATUS_CODE_SIZE            // 消息类型字段偏移量
	COMMAND_OFFSET        = MESSAGE_TYPE_OFFSET + MESSAGE_TYPE_SIZE     // 命令字段偏移量
	CLIENT_SEQ_NUM_OFFSET = COMMAND_OFFSET + COMMAND_SIZE               // 客户端序列号字段偏移量
	BODY_OFFSET           = CLIENT_SEQ_NUM_OFFSET + CLIENT_SEQ_NUM_SIZE // 消息体偏移量
)

// Header 封装消息头字段的结构体
type Header struct {
	Version      uint8  // 版本号
	StatusCode   uint8  // 状态码
	MessageType  uint16 // 消息类型
	Command      uint16 // 命令
	ClientSeqNum uint32 // 客户端序列号
}

// Message 消息结构
type Message struct {
	Header
	Body []byte // 消息体
}

// Protocol 抽象协议
type Protocol interface {
	NewParser() Parser
}

// Parser 抽象协议解析器
type Parser interface {
	Encode(Message) []byte
	Decode([]byte, uint32) (*Message, error)
}

type IMProtocol struct{}

func NewIMProtocol() Protocol {
	return &IMProtocol{}
}

func (proto *IMProtocol) NewParser() Parser {
	return &IMParser{}
}

type IMParser struct{}

// Encode 封装请求体
func (parser *IMParser) Encode(msg Message) []byte {
	/*
			# 一个-代表1个字节
		   ---- |----     |-    |-    |--    |-- |----      |Body|
		   包总长|header长度|版本号|状态码|消息类型|命令|客户端序列号|Body|
	*/
	packLen := HEADER_SIZE + len(msg.Body) + TOTAL_HEADER_SIZE
	packLenBuf := make([]byte, PACK_SIZE)
	binary.BigEndian.PutUint32(packLenBuf[:PACK_SIZE], uint32(packLen))

	buf := make([]byte, packLen)

	// header
	binary.BigEndian.PutUint16(buf[HEADER_OFFSET:], uint16(TOTAL_HEADER_SIZE))
	buf[VERSION_OFFSET] = msg.Version
	buf[STATUS_OFFSET] = msg.StatusCode
	binary.BigEndian.PutUint16(buf[MESSAGE_TYPE_OFFSET:], msg.MessageType)
	binary.BigEndian.PutUint16(buf[COMMAND_OFFSET:], msg.Command)
	binary.BigEndian.PutUint32(buf[CLIENT_SEQ_NUM_OFFSET:], msg.ClientSeqNum)

	// body
	copy(buf[HEADER_SIZE+TOTAL_HEADER_SIZE:], msg.Body)
	allBuf := append(packLenBuf, buf...)
	return allBuf
}

// Decode 解析响应体
func (parser *IMParser) Decode(data []byte, packLen uint32) (*Message, error) {
	msg := &Message{}
	msg.Version = data[VERSION_OFFSET]
	msg.StatusCode = data[STATUS_OFFSET]
	msg.MessageType = binary.BigEndian.Uint16(data[MESSAGE_TYPE_OFFSET:COMMAND_OFFSET])
	msg.Command = binary.BigEndian.Uint16(data[COMMAND_OFFSET:CLIENT_SEQ_NUM_OFFSET])
	msg.ClientSeqNum = binary.BigEndian.Uint32(data[CLIENT_SEQ_NUM_OFFSET:BODY_OFFSET])

	headerLen := binary.BigEndian.Uint16(data[HEADER_OFFSET:VERSION_OFFSET])
	if headerLen != TOTAL_HEADER_SIZE {
		return nil, fmt.Errorf("headerLen:%d, TOTAL_HEADER_SIZE:%d", headerLen, TOTAL_HEADER_SIZE)
	}
	if packLen > uint32(headerLen) {
		msg.Body = data[BODY_OFFSET:packLen]
	}
	return msg, nil
}
