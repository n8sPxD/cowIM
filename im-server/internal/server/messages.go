package server

import (
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

func (s *Server) sendMessageToBackend() error {
	// 发送到消息队列处理
	for {
		mqMsg := kafka.Message{
			Value: []byte(<-s.messages),
		}
		if err := s.svcCtx.MsgForwarder.WriteMessages(s.ctx, mqMsg); err != nil {
			logx.Error("[sendMessageToBackend] Push message to MQ failed, error: ", err)
			return err
		}
	}
}

func (s *Server) readMessageFromFrontend(id uint32) error {
	// 此处的消息为protobuf序列化后的消息
	for {
		msg, err := s.Manager.ReadMessage(id)
		if err != nil {
			// 用户断线
			if websocket.IsCloseError(
				err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				logx.Infof("[readMessageFromFrontend] User %d disconnected", id)
			} else {
				// 用户无法读取信息，不是因为用户断开了，可能是服务器的问题
				logx.Errorf(
					"[readMessageFromFrontend] Read message from user %v failed, error: %v",
					id,
					err,
				)
			}
			// 读不了用户的消息，说明后续也不能和用户进行通讯了，所以就和用户断开连接
			logx.Infof("[readMessageFromFrontend] Removing user %d from redis router status...", id)
			s.svcCtx.Redis.RemoveUserRouterStatus(s.ctx, id)
			// TODO: 这里出错要么是没读到在线状态，不影响，要么就是redis挂了，再考虑后续怎么处理
			return err
		}
		s.messages <- string(msg)
	}
}
