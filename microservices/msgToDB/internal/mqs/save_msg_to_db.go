package mqs

import (
	"context"
	"encoding/json"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/db/myMongo/models"
	"github.com/n8sPxD/cowIM/common/message/inside"
	"github.com/n8sPxD/cowIM/microservices/msgToDB/internal/svc"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type MsgToDB struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	MsgToDB *kafka.Reader
}

func NewMsgToDB(ctx context.Context, svcCtx *svc.ServiceContext) *MsgToDB {
	return &MsgToDB{
		ctx:    ctx,
		svcCtx: svcCtx,
		MsgToDB: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     svcCtx.Config.MsgToDB.Brokers,
			Topic:       svcCtx.Config.MsgToDB.Topic,
			StartOffset: kafka.LastOffset,
		}),
	}
}

func (l *MsgToDB) Close() {
	l.MsgToDB.Close()
}

func (l *MsgToDB) Start() {
	// 设置kafka起始偏移量，在初始化NewReader的时候设置没用不知道为什么，只有这里有用
	err := l.MsgToDB.SetOffset(kafka.LastOffset)
	if err != nil {
		logx.Error("[MsgForwarder.Start] Set kafka offset failed, error: ", err)
	}

	for {
		msg, err := l.MsgToDB.ReadMessage(l.ctx) // 这里的msg是kafka.Message
		if err != nil {
			logx.Error("[MsgToDB.Start] Reading msg error: ", err)
			continue
		}
		go l.Consume(msg.Value)
	}
}

func (l *MsgToDB) Consume(rawjson []byte) {
	var ins inside.Message
	err := json.Unmarshal(rawjson, &ins)
	if err != nil {
		logx.Error("[MsgToDB.Consume] Unmarshal inside.Message failed, error: ", err)
		return
	}

	// 根据类型来解析消息
	switch ins.Type {
	case constant.MESSAGE_RECORD:
		go l.saveMessageRecord(ins.Payload)
	default:
		logx.Error("[MsgToDB.Consume] Invalid type, type is: ", ins.Type)
	}
}

func (l *MsgToDB) saveMessageRecord(rawjson []byte) {
	var what models.MessageRecord
	err := json.Unmarshal(rawjson, &what)
	if err != nil {
		logx.Error("[MsgToDB.saveMessageRecord] Unmarshal message failed, error: ", err)
		return
	}
	_, err = l.svcCtx.Mongo.MessageRecord.InsertOne(l.ctx, what)
	if err != nil {
		logx.Error("[MsgToDB.saveMessageRecord] Insert message to MongoDB failed, error: ", err)
		return
	}
}
