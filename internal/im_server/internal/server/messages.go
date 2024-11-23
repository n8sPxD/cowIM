package server

import (
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"
)

func (s *Server) sendMessageToBackend() error {
	logx.Debug("[sendMessageToBackend] Sending message to MQ...")
	// 发送到消息队列处理
	for {
		mqMsg := kafka.Message{
			Value: []byte(<-s.messages),
		}
		logx.Debug("[sendMessageToBackend] Sending message: ", string(mqMsg.Value))
		if err := s.svcCtx.MsgForwarder.WriteMessages(s.ctx, mqMsg); err != nil {
			logx.Error("[sendMessageToBackend] Push message to MQ failed, error: ", err)
			return err
		}
		logx.Debug("[sendMessageToBackend] Send message success")
	}
}

func (s *Server) readMessageFromFrontend(id uint32) error {
	for {
		msg, err := s.Manager.ReadMessage(id)
		if err != nil {
			// 用户断线
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				logx.Infof("[readMessageFromFrontend] User %d disconnected", id)
			} else {
				// 用户无法读取信息，不是因为用户断开了，可能是服务器的问题
				logx.Errorf("[readMessageFromFrontend] Read message from user %v failed, error: %v", id, err)
			}
			return err
		}
		// 判定消息类型
		// TODO: 这里的消息类型可以放到MsgForwarder中去判断
		if isHeartbeat(msg) {
			go s.checkHeartBeat(id)
		} else if isAck(msg) {
			go s.acceptAckMessage(msg)
		} else {
			s.messages <- string(msg)
		}
	}
}

func isHeartbeat(message []byte) bool {
	return string(message) == "ping"
}

func isAck(message []byte) bool {
	return strings.HasPrefix(string(message), "ack")
}

// 心跳检查
func (s *Server) checkHeartBeat(id uint32) {
	//if _, err := s.svcCtx.Redis.UpdateUserRouterStatus(s.ctx, id, s.svcCtx.Config.WorkID); err != nil {
	//	logx.Error("[checkHeartBeat] Update router status to redis failed, error: ", err)
	//} else {
	//	logx.Info("HeartBeat from User ", id)
	//}
}

// 接受ACK消息
func (s *Server) acceptAckMessage(message []byte) {
	// 先从string把数据分出来
	// 0: ack 1: userID 2: messageID
	parts := strings.Split(string(message), "_")
	tmpID, _ := strconv.Atoi(parts[1])
	userID := uint32(tmpID)
	messageID := parts[2]

	s.Manager.GetAckHandler().ConfirmAck(Ack{
		To:        userID,
		MessageID: messageID,
	})

	logx.Debugf("[acceptAckMessage] Confirm ack message to user %d with messageID \"%s\"", userID, messageID)
}
