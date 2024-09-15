// Code generated by goctl. DO NOT EDIT.
// Source: im-rpc.proto

package imrpcclient

import (
	"context"

	"github.com/n8sPxD/cowIM/imrpc/imrpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	SendMessageRequest  = imrpc.SendMessageRequest
	SendMessageResponse = imrpc.SendMessageResponse

	ImRPC interface {
		SendSingleMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
	}

	defaultImRPC struct {
		cli zrpc.Client
	}
)

func NewImRPC(cli zrpc.Client) ImRPC {
	return &defaultImRPC{
		cli: cli,
	}
}

func (m *defaultImRPC) SendSingleMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	client := imrpc.NewImRPCClient(m.cli.Conn())
	return client.SendSingleMessage(ctx, in, opts...)
}
