package logic

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/n8sPxD/cowIM/internal/common/constant"
	models2 "github.com/n8sPxD/cowIM/internal/common/dao/myMongo/models"
	"github.com/n8sPxD/cowIM/internal/common/message/front"
	"github.com/n8sPxD/cowIM/internal/common/message/inside"
	"github.com/segmentio/kafka-go"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/zeromicro/go-zero/core/logx"
)

func (l *MsgForwarder) sendTimelineToDB(msg *front.Message, now time.Time) {
	id, _ := strconv.ParseInt(msg.Id, 10, 64)
	syncMsg := models2.MessageSync{
		ID:        id,
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
		senderTL := models2.UserTimeline{
			ID:         idgen.NextId(),
			ReceiverID: msg.To,
			SenderID:   msg.From,
			Type:       constant.SINGLE_CHAT,
			Message:    syncMsg,
			Timestamp:  now,
		}
		senderTLByte, err := json.Marshal(senderTL)
		if err != nil {
			logx.Error("[sendTimelineToDB] Json marshal failed, error: ", err)
			return
		}

		var packMsg inside.MessageToDB
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
		err = l.MsgDBSaver.WriteMessages(l.ctx, mqMsg)
		if err != nil {
			logx.Error("[sendTimelineToDB] Push message to DBSaver MQ failed, error: ", err)
		}

	case constant.GROUP_CHAT:
		// 先从MySQL查找群聊成员
		members, err := l.svcCtx.MySQL.GetGroupMemberIDs(l.ctx, uint(*msg.Group))
		if err != nil {
			logx.Error("[sendTimelineToDB] GetGroupMembers failed, error: ", err)
			return
		}

		// 然后封装所有群员的消息
		var pack inside.MessageToDB
		pack.Type = constant.USER_TIMELINE
		for _, member := range members {
			current := models2.UserTimeline{
				ID:         idgen.NextId(),
				ReceiverID: uint32(member),
				SenderID:   msg.From,
				Type:       constant.GROUP_CHAT,
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
		err = l.MsgDBSaver.WriteMessages(l.ctx, mqMsg)
		if err != nil {
			logx.Error("[sendTimelineToDB] Push message to DBSaver MQ failed, error: ", err)
		}

	case constant.BIG_GROUP_CHAT:
		// 大规模群聊，读扩散存库

		// TODO:
		//senderTLByte, err := json.Marshal(senderTL)
		//if err != nil {
		//	logx.Error("[sendTimelineToDB] Json marshal failed, error: ", err)
		//	return
		//}

		var packMsg inside.MessageToDB
		packMsg.Type = constant.USER_TIMELINE
		packMsg.Payload = append(
			packMsg.Payload,
			// TODO: senderTLByte,
		)
		packMsgByte, err := json.Marshal(packMsg)
		if err != nil {
			logx.Error("[sendTimelineToDB] Json marshal failed, error: ", err)
			return
		}
		mqMsg := kafka.Message{
			Value: packMsgByte,
		}
		err = l.MsgDBSaver.WriteMessages(l.ctx, mqMsg)
		if err != nil {
			logx.Error("[sendTimelineToDB] Push message to DBSaver MQ failed, error: ", err)
		}
	case constant.SYSTEM_INFO:
	default:
		logx.Error("[sendTimelineToDB] Invalid message type, type is: ", msg.Type)
	}
}

// 封装消息，发送到存库服务中进行存库
// 对于单聊消息，正常存储
// 对于群组消息， SenderID对应 front.Message 中的from
// ReceiverID 对应 front.Message 中的 to，group也可以
func (l *MsgForwarder) sendRecordMsgToDB(msg *front.Message, now time.Time) {
	id, _ := strconv.ParseInt(msg.Id, 10, 64)
	recordMsg := models2.MessageRecord{
		ID:         id,
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

	var packMsg inside.MessageToDB
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
	err = l.MsgDBSaver.WriteMessages(l.ctx, mqMsg)
	if err != nil {
		logx.Error("[sendRecordMsgToDB] Push message to DBSaver MQ failed, error: ", err)
	}
}
