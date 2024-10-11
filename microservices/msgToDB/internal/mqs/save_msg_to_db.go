package mqs

import (
	"context"
	"encoding/json"
	"time"

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
			GroupID:     "msg-fwd",
			StartOffset: kafka.LastOffset,
			// TODO: 入库服务对延迟要求不高，可以适当放宽参数条件
			MinBytes:       1,                      // 最小拉取字节数
			MaxBytes:       10e3,                   // 最大拉取字节数（10KB）
			MaxWait:        100 * time.Millisecond, // 最大等待时间
			CommitInterval: 500 * time.Millisecond, // 提交间隔
		}),
	}
}

func (l *MsgToDB) Close() {
	l.MsgToDB.Close()
}

func (l *MsgToDB) Start() {
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
	var ins inside.MessageToDB
	err := json.Unmarshal(rawjson, &ins)
	if err != nil {
		logx.Error("[MsgToDB.Consume] Unmarshal inside.MessageToDB failed, error: ", err)
		return
	}

	// 根据类型来解析消息
	switch ins.Type {
	case constant.MESSAGE_RECORD:
		go l.saveMessageRecord(ins.Payload[0]) // recordMessage 一次只会存一条
	case constant.USER_TIMELINE:
		go l.saveUserTimeline(ins.Payload)
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

func (l *MsgToDB) saveUserTimeline(rawjson []json.RawMessage) {
	var huh []interface{}
	for _, bytes := range rawjson {
		var tmp models.UserTimeline
		err := json.Unmarshal(bytes, &tmp)
		if err != nil {
			logx.Error("[MsgToDB.saveMessageRecord] Unmarshal message failed, error: ", err)
			return
		}
		huh = append(huh, tmp)
	}

	_, err := l.svcCtx.Mongo.TimeLine.InsertMany(l.ctx, huh)
	if err != nil {
		logx.Error("[MsgToDB.saveMessageRecord] Insert message to MongoDB failed, error: ", err)
		return
	}
}
