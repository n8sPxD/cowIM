package mqs

import (
	"context"
	"time"

	"github.com/n8sPxD/cowIM/common/constant"
	"github.com/n8sPxD/cowIM/common/libnet"
	"github.com/n8sPxD/cowIM/common/socket"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/config"
	"github.com/n8sPxD/cowIM/gateway/tcp_gateway/svc"
	"github.com/n8sPxD/cowIM/imrpc/imrpc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"google.golang.org/protobuf/proto"
)

type Chat struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	server *socket.Server
}

func NewChat(ctx context.Context, svcCtx *svc.ServiceContext, server *socket.Server) *Chat {
	return &Chat{
		ctx:    ctx,
		svcCtx: svcCtx,
		server: server,
	}
}

func (l *Chat) Consume(ctx context.Context, key, val string) error {
	var msg imrpc.SendMessageRequest
	if err := proto.Unmarshal([]byte(val), &msg); err != nil {
		logx.Error("Consume message unmarshal error, err: ", err)
		return err
	}
	// 传入参数中没有制定接收用户，则为群聊
	if msg.RecvUser == constant.NONE_SINGLE {
		// TODO: 群聊
		l.groupChat()
	} else if msg.RecvGroup == constant.NONE_GROUP {
		logx.Info("single chat")
		// TODO: 单聊
		if err := l.singleChat(&msg); err != nil {
			logx.Error("Consume singleChat error, err: ", err)
			return err
		}
	}
	return nil

}

func (l *Chat) groupChat() {}

func (l *Chat) singleChat(msg *imrpc.SendMessageRequest) error {
	recvUser := msg.RecvUser
	if s, online := l.server.Manager.Sessions[recvUser]; online {
		logx.Infof("user %s online", recvUser)
		// 用户在线，直接调用Send发送
		// TODO: 用户那边指定消息header，通过改写rpc的request
		realMsg := libnet.Message{
			Header: libnet.Header{
				Version:      1,
				StatusCode:   200,
				MessageType:  constant.COMMON_MSG,
				Command:      0,
				ClientSeqNum: 0, // TODO: 使用全局唯一递增ID发送器生成
			},
			Body: []byte(msg.Content),
		}
		logx.Infof("Sending message \"%v\" to user \"%s\"", realMsg, s.User())
		if err := s.Send(realMsg); err != nil {
			return err
		}
	}
	return nil
}

func Consumers(c config.Config, ctx context.Context, svcContext *svc.ServiceContext, server *socket.Server) []service.Service {
	return []service.Service{
		kq.MustNewQueue(c.KqConsumerConf, NewChat(ctx, svcContext, server), kq.WithCommitInterval(time.Millisecond*500)),
	}
}
