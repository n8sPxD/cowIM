package servicehub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/n8sPxD/cowIM/pkg/loadbalancer"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type DiscoveryHub struct {
	client       *clientv3.Client
	loadBalancer loadbalancer.LoadBalancer
}

var (
	discoveryHub  *DiscoveryHub
	discoveryOnce sync.Once
)

// NewDiscoveryHub 单例模式创建一个DiscoveryHub
func NewDiscoveryHub(etcdServers []string, heartbeatFrequency int64) *DiscoveryHub {
	if discoveryHub == nil {
		discoveryOnce.Do(func() {
			if client, err := clientv3.New(
				clientv3.Config{
					Endpoints:   etcdServers,
					DialTimeout: 3 * time.Second,
				}); err != nil {
				logx.Error("[GetDiscoveryHub] Connect to etcd failed, error: ", err)
			} else {
				discoveryHub = &DiscoveryHub{
					client:       client,
					loadBalancer: loadbalancer.NewLoadBalancer(loadbalancer.RoundRobinBalancer),
				}
			}
		})
	}
	return discoveryHub
}

func (hub *DiscoveryHub) GetServiceEndpoints(ctx context.Context, service string) []string {
	prefix := fmt.Sprintf("%s/", service)
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

// GetServiceEndpoint 通过负载均衡算法获取一个服务地址
func (hub *DiscoveryHub) GetServiceEndpoint(ctx context.Context, service string) string {
	return hub.loadBalancer.Take(hub.GetServiceEndpoints(ctx, service))
}

func (hub *DiscoveryHub) Close() {
	hub.client.Close()
}
