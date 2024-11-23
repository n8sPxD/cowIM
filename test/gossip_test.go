package test

import (
	"context"
	"fmt"
	"github.com/n8sPxD/cowIM/internal/msg_forward/gossip"
	"github.com/n8sPxD/cowIM/internal/msg_forward/gossip/gossippb"
	"github.com/n8sPxD/cowIM/pkg/servicehub"
	"github.com/n8sPxD/cowIM/pkg/utils"
	"testing"
	"time"
)

func TestGossip(t *testing.T) {
	discov := servicehub.NewDiscoveryHub([]string{"127.0.0.1:2379"})
	regist := servicehub.NewRegisterHub([]string{"127.0.0.1:2379"}, 5)

	server1 := gossip.MustNewServer(discov, regist, 6666, 1, 1, 2, 3)
	server2 := gossip.MustNewServer(discov, regist, 6667, 2, 1, 2, 3)

	go server1.Start()
	go server2.Start()

	ip, err := utils.GetLocalIP()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	data := make([]*gossippb.Data, 0)
	temp := gossippb.Data{
		Key:       233,
		Value:     6,
		Version:   1,
		Timestamp: time.Now().UnixMilli(),
	}
	data = append(data, &temp)

	time.Sleep(6 * time.Second)

	client1 := gossip.NewClient(fmt.Sprintf("%s:%d", ip, 6666))
	if _, err := client1.RemoteUpdate(context.Background(), &gossippb.RemoteRequest{Data: data}); err != nil {
		t.Error(err)
		t.Fail()
	}

	// 等待同步
	time.Sleep(6 * time.Second)

	if data, ok := server1.Node.Data[233]; !ok {
		t.Fail()
	} else {
		t.Logf("value: %d, version: %d, timestamp: %d", data.Value, data.Version, data.Timestamp)
	}

	if data, ok := server2.Node.Data[233]; !ok {
		t.Fail()
	} else {
		t.Logf("value: %d, version: %d, timestamp: %d", data.Value, data.Version, data.Timestamp)
	}
}
