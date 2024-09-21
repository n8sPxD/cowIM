package mqs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/db/myMongo/models"
	"github.com/n8sPxD/cowIM/common/message/front"
	"github.com/n8sPxD/cowIM/common/message/inside"
	"github.com/n8sPxD/cowIM/microservices/msgForward/internal/svc"
	"github.com/segmentio/kafka-go"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type MsgForwarder struct {
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	MsgForwarder *kafka.Reader
}

func NewMsgForwarder(ctx context.Context, svcCtx *svc.ServiceContext) *MsgForwarder {
	return &MsgForwarder{
		ctx:    ctx,
		svcCtx: svcCtx,
		MsgForwarder: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     svcCtx.Config.MsgForwarder.Brokers,
			Topic:       svcCtx.Config.MsgForwarder.Topic,
			StartOffset: kafka.LastOffset,
		}),
	}
}

func (l *MsgForwarder) Start() {
	// 初始化id生成器
	options := idgen.NewIdGeneratorOptions(l.svcCtx.Config.WorkID)
	idgen.SetIdGenerator(options)

	for {
		msg, err := l.MsgForwarder.ReadMessage(l.ctx) // 这里的msg是kafka.Message
		if err != nil {
			logx.Error("[MsgForwarder.Start] Reading msgForward error: ", err)
			continue
		}
		logx.Infof(
			"message at offset %d: %s = %s\n",
			msg.Offset,
			string(msg.Key),
			string(msg.Value),
		)
		go l.Consume(msg.Value, time.Now())
	}
}

// Consume 接收从 Websocket Server的消息，处理后再进行转发
func (l *MsgForwarder) Consume(protobuf []byte, current time.Time) {
	// 传过来的消息是序列化过的，先反序列化
	var msg front.Message
	err := proto.Unmarshal(protobuf, &msg)
	if err != nil {
		logx.Error("[MsgForwarder.Consume] Unmarshal msgForward failed, error: ", err)
		return
	}

	// 异步存库
	go l.sendRecordMsgToDB(&msg, current)

	// 进行基于消息类型的消息处理
	switch msg.Type {
	case constant.SINGLE_CHAT:
		l.singleChat(&msg, protobuf)
	case constant.GROUP_CHAT:
		l.groupChat(&msg)
	case constant.BIG_GROUP_CHAT:
		l.bigGroupChat(&msg)
	default:
		logx.Error("[MsgForwarder.Consume] Wrong msgForward type, Type is: ", msg.Type)
	}
}

func (l *MsgForwarder) singleChat(msg *front.Message, protobuf []byte) {
	// 查询Redis中路由信息
	status, err := l.svcCtx.Redis.GetUserRouterStatus(l.ctx, msg.To)
	if err != nil {
		logx.Error("[MsgForwarder.singleChat] Get router status from redis failed, error: ", err)
		return
	}
	if status == nil {
		// 没找到当前用户的路由信息，说明没上线
		// TODO: 塞timeline里
		return
	}
	// 转发消息到指定的websocket-server
	// 先确定Topic
	workID := status.WorkID
	l.svcCtx.MsgSender.Topic = fmt.Sprintf("websocket-server-%d", workID)
	// 封装消息
	mqMsg := kafka.Message{
		Value: protobuf,
	}
	err = l.svcCtx.MsgSender.WriteMessages(l.ctx, mqMsg)
	if err != nil {
		logx.Error(
			"[MsgForwarder.singleChat] Push msgForward to Websocket-server MQ failed, error: ",
			err,
		)
		return
	}
}

func (l *MsgForwarder) groupChat(msg *front.Message) {
	// TODO: 完善逻辑
}

func (l *MsgForwarder) bigGroupChat(msg *front.Message) {
	// TODO: 完善逻辑
}

// 封装消息，发送到存库服务中进行存库
func (l *MsgForwarder) sendRecordMsgToDB(msg *front.Message, current time.Time) {
	recordMsg := models.MessageRecord{
		ID:         idgen.NextId(),
		SenderID:   msg.From,
		Type:       uint8(msg.Type),
		ReceiverID: msg.To,
		MsgType:    uint8(msg.MsgType),
		Content:    msg.Content,
		Timestamp:  current,
	}
	if msg.Extend != nil {
		recordMsg.Extend = *msg.Extend
	}

	rawRecordMsg, err := json.Marshal(recordMsg)
	if err != nil {
		logx.Error("[MsgForwarder.Consume] Marshal MessageRecord failed, error: ", err)
		return
	}

	packMsg := inside.Message{
		Type:    constant.MESSAGE_RECORD,
		Payload: rawRecordMsg,
	}
	rawPackMsg, err := json.Marshal(packMsg)
	if err != nil {
		logx.Error("[MsgForwarder.Consume] Marshal PackMessageRecord failed, error: ", err)
		return
	}

	mqMsg := kafka.Message{
		Value: rawPackMsg,
	}
	err = l.svcCtx.MsgDBSaver.WriteMessages(l.ctx, mqMsg)
	if err != nil {
		logx.Error("[MsgForwarder.Consume] Push message to DBSaver MQ failed, error: ", err)
	}
}
