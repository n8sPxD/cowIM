package gossip

import (
	"context"
	"github.com/n8sPxD/cowIM/pkg/servicehub"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"sync"
	"time"
)

type NodeState struct {
	NodeID     string
	Status     map[string]string // 路由状态以及其他数据
	Version    int64             // 状态版本号
	LastUpdate time.Time         // 最后更新时间戳
}

type GossipNode struct {
	mu          sync.Mutex
	NodeID      string
	LocalState  NodeState            // 本地状态
	NeighborMap map[string]NodeState // 邻居状态
}

func NewGossipNode(nodeID string) *GossipNode {
	return &GossipNode{
		NodeID: nodeID,
		LocalState: NodeState{
			NodeID:     nodeID,
			Status:     make(map[string]string),
			Version:    0,
			LastUpdate: time.Now(),
		},
		NeighborMap: make(map[string]NodeState),
	}
}

type GossipMessage struct {
	SourceNode string               // 消息来源
	StateMap   map[string]NodeState // 包含状态数据
}

func (node *GossipNode) StartGossip(neighbors []string) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		if len(neighbors) == 0 {
			continue
		}
		random := neighbors[rand.Intn(len(neighbors))]
		client, err := NewGossipClient(random)
		if err != nil {
			logx.Errorf("Failed to connect to neighbor %s: %v", random, err)
			continue
		}
		defer client.conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		node.mu.Lock()
		message := GossipMessage{
			SourceNode: node.NodeID,
			StateMap:   map[string]NodeState{node.NodeID: node.LocalState},
		}
		for k, v := range node.NeighborMap {
			message.StateMap[k] = v
		}
		node.mu.Unlock()

		if err := client.PropagateGossip(ctx, message); err != nil {
			logx.Errorf("Failed to propagate gossip to %s: %v", random, err)
		}
	}
}

func (node *GossipNode) HandleGossipMessage(message GossipMessage) {
	node.mu.Lock()
	defer node.mu.Unlock()

	for id, received := range message.StateMap {
		if id == node.NodeID {
			continue
		}
		// 如果收到的版本号更高或者时间更新，更新本地邻居状态
		local, exists := node.NeighborMap[id]
		if !exists || received.Version > local.Version || received.LastUpdate.After(local.LastUpdate) {
			node.NeighborMap[id] = received
		}
	}
}

func (node *GossipNode) UpdateLocalState(key, value string) {
	node.mu.Lock()
	defer node.mu.Unlock()

	node.LocalState.Status[key] = value
	node.LocalState.Version++
	node.LocalState.LastUpdate = time.Now()
}

func (node *GossipNode) PushUpdate(neighbors []string, sendMessage func(string, GossipMessage)) {
	node.mu.Lock()
	message := GossipMessage{
		SourceNode: node.NodeID,
		StateMap:   map[string]NodeState{node.NodeID: node.LocalState},
	}
	node.mu.Unlock()

	for _, neighbor := range neighbors {
		sendMessage(neighbor, message)
	}
}

func (node *GossipNode) PullState(neighbor string, fetchSize func(string) map[string]NodeState) {
	received := fetchSize(neighbor)

	node.mu.Lock()
	defer node.mu.Unlock()

	for id, state := range received {
		if id == node.NodeID {
			continue
		}
		local, exists := node.NeighborMap[id]
		if !exists || state.Version > local.Version || state.LastUpdate.After(local.LastUpdate) {
			node.NeighborMap[id] = state
		}
	}
}

func (node *GossipNode) DiscoverNeighbors(ctx context.Context, hub *servicehub.DiscoveryHub) []string {
	return hub.GetServiceEndpoints(ctx, "message-forward")
}
