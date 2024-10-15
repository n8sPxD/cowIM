package server

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
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
		if string(msg) == "ping" {
			go s.checkHeartBeat(id)
		} else {
			s.messages <- string(msg)
		}
	}
}

// 心跳检查
func (s *Server) checkHeartBeat(id uint32) {
	if err := s.svcCtx.Redis.UpdateUserRouterStatus(s.ctx, id, s.svcCtx.Config.WorkID, time.Now()); err != nil {
		logx.Error("[checkHeartBeat] Update router status to redis failed, error: ", err)
	}
	logx.Info("HeartBeat from User ", id)
}
