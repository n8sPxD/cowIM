package mqs

import (
	"encoding/json"
	"time"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/db/myMongo/models"
	"github.com/n8sPxD/cowIM/common/message/front"
	"github.com/n8sPxD/cowIM/common/message/inside"
	"github.com/segmentio/kafka-go"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/zeromicro/go-zero/core/logx"
)

func (l *MsgForwarder) sendTimelineToDB(msg *front.Message, now time.Time) {
	syncMsg := models.MessageSync{
		ID:        idgen.NextId(),
		MsgType:   uint8(msg.MsgType),
		Content:   msg.Content,
		Timestamp: now,
	}
	if msg.Extend != nil {
		syncMsg.Extend = *msg.Extend
	}

	// 分别处理群聊和单聊
	switch msg.Type {
	case constant.SINGLE_CHAT:
		senderTL := models.UserTimeline{
			ID:         idgen.NextId(),
			ReceiverID: msg.To,
			SenderID:   msg.From,
			// GroupID:   0,  	// 到下面去判断
			Message:   syncMsg,
			Timestamp: now,
		}
		senderTLByte, err := json.Marshal(senderTL)
		if err != nil {
			logx.Error("[sendTimelineToDB] Json marshal failed, error: ", err)
			return
		}

		var packMsg inside.Message
		packMsg.Type = constant.USER_TIMELINE
		packMsg.Payload = append(
			packMsg.Payload,
			senderTLByte,
		)
		packMsgByte, err := json.Marshal(packMsg)
		if err != nil {
			logx.Error("[sendTimelineToDB] Json marshal failed, error: ", err)
			return
		}
		mqMsg := kafka.Message{
			Value: packMsgByte,
		}
		err = l.svcCtx.MsgDBSaver.WriteMessages(l.ctx, mqMsg)
		if err != nil {
			logx.Error("[sendTimelineToDB] Push message to DBSaver MQ failed, error: ", err)
		}

	case constant.GROUP_CHAT:
		// 先从MySQL查找群聊成员
		members, err := l.svcCtx.MySQL.GetGroupMembers(l.ctx, uint(msg.To))
		if err != nil {
			logx.Error("[sendTimelineToDB] GetGroupMembers failed, error: ", err)
			return
		}

		// 然后封装所有群员的消息
		var pack inside.Message
		pack.Type = constant.USER_TIMELINE
		for _, member := range members {
			current := models.UserTimeline{
				ID:         idgen.NextId(),
				ReceiverID: member,
				SenderID:   msg.From,
				GroupID:    msg.To,
				Message:    syncMsg,
				Timestamp:  now,
			}
			// json序列化
			currentByte, err := json.Marshal(current)
			if err != nil {
				logx.Error("[sendTimelineToDB] Json marshal failed, error: ", err)
				return
			}
			// 封装到打包中
			pack.Payload = append(pack.Payload, currentByte)
		}

		packByte, err := json.Marshal(pack)
		if err != nil {
			logx.Error("[sendTimelineToDB] Json marshal failed, error: ", err)
			return
		}

		mqMsg := kafka.Message{
			Value: packByte,
		}
		err = l.svcCtx.MsgDBSaver.WriteMessages(l.ctx, mqMsg)
		if err != nil {
			logx.Error("[sendTimelineToDB] Push message to DBSaver MQ failed, error: ", err)
		}

	case constant.BIG_GROUP_CHAT:
	case constant.SYSTEM_INFO:
	default:
		logx.Error("[sendTimelineToDB] Invalid message type, type is: ", msg.Type)
	}
}

// 封装消息，发送到存库服务中进行存库
func (l *MsgForwarder) sendRecordMsgToDB(msg *front.Message, now time.Time) {
	recordMsg := models.MessageRecord{
		ID:         idgen.NextId(),
		SenderID:   msg.From,
		Type:       uint8(msg.Type),
		ReceiverID: msg.To,
		MsgType:    uint8(msg.MsgType),
		Content:    msg.Content,
		Timestamp:  now,
	}
	if msg.Extend != nil {
		recordMsg.Extend = *msg.Extend
	}

	recordMsgByte, err := json.Marshal(recordMsg)
	if err != nil {
		logx.Error("[sendRecordMsgToDB] Marshal MessageRecord failed, error: ", err)
		return
	}

	var packMsg inside.Message
	packMsg.Type = constant.MESSAGE_RECORD
	packMsg.Payload = append(packMsg.Payload, recordMsgByte)

	rawPackMsg, err := json.Marshal(packMsg)
	if err != nil {
		logx.Error("[sendRecordMsgToDB] Marshal PackMessageRecord failed, error: ", err)
		return
	}

	mqMsg := kafka.Message{
		Value: rawPackMsg,
	}
	err = l.svcCtx.MsgDBSaver.WriteMessages(l.ctx, mqMsg)
	if err != nil {
		logx.Error("[sendRecordMsgToDB] Push message to DBSaver MQ failed, error: ", err)
	}
}
