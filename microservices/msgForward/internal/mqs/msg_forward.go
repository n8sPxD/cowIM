package mqs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/message/front"
	"github.com/n8sPxD/cowIM/common/message/inside"
	"github.com/n8sPxD/cowIM/common/utils"
	"github.com/n8sPxD/cowIM/microservices/msgForward/internal/svc"
	"github.com/segmentio/kafka-go"
	"github.com/yitter/idgenerator-go/idgen"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
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
			Brokers:        svcCtx.Config.MsgForwarder.Brokers,
			Topic:          svcCtx.Config.MsgForwarder.Topic,
			GroupID:        "msg-fwd",
			StartOffset:    kafka.LastOffset,
			MinBytes:       1,                      // 最小拉取字节数
			MaxBytes:       10e3,                   // 最大拉取字节数（10KB）
			MaxWait:        100 * time.Millisecond, // 最大等待时间
			CommitInterval: 500 * time.Millisecond, // 提交间隔
		}),
	}
}

func (l *MsgForwarder) Close() {
	l.MsgForwarder.Close()
}

func (l *MsgForwarder) Start() {
	// 初始化id生成器
	options := idgen.NewIdGeneratorOptions(l.svcCtx.Config.WorkID)
	idgen.SetIdGenerator(options)

	for {
		msg, err := l.MsgForwarder.ReadMessage(l.ctx) // 这里的msg是kafka.Message
		if err != nil {
			logx.Error("[MsgForwarder.Start] Reading msgForward error: ", err)
			break
		}
		logx.Debugf(
			"[MsgForwarder.Start] Message at partition %d offset %d: %s\n",
			msg.Partition,
			msg.Offset,
			string(msg.Value),
		)
		go l.Consume(msg.Value, time.Now())
	}
}

// Consume 接收从 Websocket Server的消息，处理后再进行转发
func (l *MsgForwarder) Consume(protobuf []byte, now time.Time) {
	// 传过来的消息是序列化过的，先反序列化
	var msg front.Message
	err := proto.Unmarshal(protobuf, &msg)
	if err != nil {
		logx.Error("[MsgForwarder.Consume] Unmarshal msgForward failed, error: ", err)
		return
	}

	// 异步存库
	go l.sendRecordMsgToDB(&msg, now) // 漫游库
	go l.sendTimelineToDB(&msg, now)  // timeline+同步库

	// 进行基于消息类型的消息处理
	switch msg.Type {
	case constant.SINGLE_CHAT:
		go l.singleChat(&msg, protobuf)
	case constant.GROUP_CHAT:
		go l.groupChat(&msg)
	case constant.BIG_GROUP_CHAT:
		go l.bigGroupChat(&msg)
	default:
		logx.Error("[MsgForwarder.Consume] Wrong msgForward type, Type is: ", msg.Type)
	}
}

// 单聊处理
// 需要传入原消息是因为单聊中的from和to无需更改,直接原样转发即可
func (l *MsgForwarder) singleChat(msg *front.Message, protobuf []byte) {
	// 查询Redis中路由信息
	status, err := l.svcCtx.Redis.GetUserRouterStatus(l.ctx, msg.To)
	if errors.Is(err, redis.Nil) {
		// 没找到当前用户的路由信息，说明没上线
		// 之前已经存过timeline了，所以不需要做任何处理
		return
	}
	if err != nil {
		logx.Error("[MsgForwarder.singleChat] Get router status from redis failed, error: ", err)
		return
	}

	// 转发消息到指定的websocket-server
	// 先确定Topic
	workID := status.WorkID
	l.svcCtx.MsgSender.Topic = fmt.Sprintf("websocket-server-%d", workID)

	// 封装inside.Message
	wsmsg := inside.Message{
		To:       msg.To,
		Protobuf: protobuf,
	}
	wsmsgByte, err := proto.Marshal(&wsmsg)
	if err != nil {
		logx.Errorf(utils.FmtFuncName(), " Marshal message failed, error: ", err)
		return
	}
	// 封装mq消息
	mqMsg := kafka.Message{
		Value: wsmsgByte,
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

// 群聊处理
func (l *MsgForwarder) groupChat(msg *front.Message) {
	// 先获取群里所有成员
	members, err := l.svcCtx.MySQL.GetGroupMembers(l.ctx, uint(msg.To))
	if err != nil {
		logx.Error(utils.FmtFuncName(), " Get group members from mysql failed, error: ", err)
		return
	}
	// TODO: 可以优化，先处理所有消息，然后把对应服务器的消息以切片形式发送，避免重复调用WriteMessages
	for _, member := range members {
		receiver := member.UserID
		// 先查用户在不在线
		status, err := l.svcCtx.Redis.GetUserRouterStatus(l.ctx, uint32(receiver))
		if errors.Is(err, redis.Nil) {
			// 不在线，直接跳过
			continue
		} else if err != nil {
			// redis出问题了
			logx.Error(utils.FmtFuncName(), " Get router status from redis failed, error: ", err)
			continue
		}
		// 用户在线，发消息
		// 确定Topic
		workID := status.WorkID
		l.svcCtx.MsgSender.Topic = fmt.Sprintf("websocket-server-%d", workID)
		mqmsg := kafka.Message{
			// TODO: 塞序列化好的消息
		}
		// 发消息
		if err := l.svcCtx.MsgSender.WriteMessages(l.ctx, mqmsg); err != nil {
			logx.Error(utils.FmtFuncName(), " Push msgForward to Websocket-server MQ failed, error: ", err)
			continue
		}
	}
}

func (l *MsgForwarder) bigGroupChat(msg *front.Message) {
	// TODO: 完善逻辑
}
