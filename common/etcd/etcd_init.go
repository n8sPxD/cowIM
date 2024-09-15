package etcd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func MustNewEtcd(add string) *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{add},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return cli
}
