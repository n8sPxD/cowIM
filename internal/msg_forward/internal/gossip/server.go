package gossip

import (
	"context"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/gossip/rpc"
	"time"
)

type GossipServiceServer struct {
	node *GossipNode
}

func (s *GossipServiceServer) PropagateGossip(ctx context.Context, req *rpc.GossipMessage) (*rpc.GossipResponse, error) {
	message := GossipMessage{
		SourceNode: req.SourceNode,
		StateMap:   make(map[string]NodeState),
	}
	for k, v := range req.StateMap {
		message.StateMap[k] = NodeState{
			NodeID:     v.NodeId,
			Status:     v.Status,
			Version:    v.Version,
			LastUpdate: parseTime(v.LastUpdate),
		}
	}
	s.node.HandleGossipMessage(message)
	return &rpc.GossipResponse{Ack: "success"}, nil
}

func (s *GossipServiceServer) PullState(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	s.node.mu.Lock()
	defer s.node.mu.Unlock()

	response := &rpc.PullResponse{StateMap: make(map[string]*rpc.NodeState)}
	for k, v := range s.node.NeighborMap {
		response.StateMap[k] = &rpc.NodeState{
			NodeId:     v.NodeID,
			Status:     v.Status,
			Version:    v.Version,
			LastUpdate: v.LastUpdate.Format(time.RFC3339),
		}
	}
	return response, nil
}
