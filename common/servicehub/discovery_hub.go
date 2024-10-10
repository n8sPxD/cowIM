package servicehub

import (
	"context"
	"fmt"

	"github.com/n8sPxD/cowIM/common/lb"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type DiscoveryHub struct {
	client       clientv3.Client
	loadBalancer lb.LoadBalancer
}

func (hub *DiscoveryHub) GetServiceEndpoints(ctx context.Context) []string {
	prefix := fmt.Sprintf("%s/", KEY_PREFIX)
	if resp, err := hub.client.Get(ctx, prefix, clientv3.WithPrefix()); err != nil {
		logx.Error("[GetServiceEndpoints] Get service from etcd failed, error: ", err)
		return nil
	} else {
		endpoints := make([]string, 0, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			endpoints = append(endpoints, string(kv.Value))
		}
		return endpoints
	}
}

func (hub *DiscoveryHub) GetServiceEndpoint(ctx context.Context) string {
	return hub.loadBalancer.Take(hub.GetServiceEndpoints(ctx))
}

func (hub *DiscoveryHub) Close() {
	hub.client.Close()
}
