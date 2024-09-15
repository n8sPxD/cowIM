// etcd/register.go
// 服务注册与发现

package etcd

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
)

// DeliveryAddress 服务注册
func DeliveryAddress(etcdAddr string, serviceName string, addr string) {
	list := strings.Split(addr, ":")
	if len(list) != 2 {
		logx.Errorf("ip error %s", addr)
		return
	}
	if list[0] == "0.0.0.0" {
		ip := netx.InternalIp()
		strings.ReplaceAll(addr, "0.0.0.0", ip)
	}

	client := MustNewEtcd(etcdAddr)
	_, err := client.Put(context.Background(), serviceName, addr)
	if err != nil {
		logx.Errorf("etcd 连接失败 %s", err.Error())
		return
	}
	logx.Infof("服务上送成功 %s %s", serviceName, addr)
}

// GetServiceAddr 服务发现
func GetServiceAddr(etcdAddr string, serviceName string) (addr string) {
	client := MustNewEtcd(etcdAddr)
	res, err := client.Get(context.Background(), serviceName)
	if err == nil && len(res.Kvs) > 0 {
		return string(res.Kvs[0].Value)
	}
	return ""
}
