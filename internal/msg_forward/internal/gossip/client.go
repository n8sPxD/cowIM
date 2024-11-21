package gossip

import (
	"context"
	"github.com/n8sPxD/cowIM/internal/msg_forward/internal/gossip/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type GossipClient struct {
	conn   *grpc.ClientConn
	client rpc.GossipServiceClient
}

func NewGossipClient(target string) (*GossipClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GossipClient{
		conn:   conn,
		client: rpc.NewGossipServiceClient(conn),
	}, nil
}

func (gc *GossipClient) PropagateGossip(ctx context.Context, message GossipMessage) error {
	stateMap := make(map[string]*rpc.NodeState)
	for k, v := range message.StateMap {
		stateMap[k] = &rpc.NodeState{
			NodeId:     v.NodeID,
			Status:     v.Status,
			Version:    v.Version,
			LastUpdate: v.LastUpdate.Format(time.RFC3339),
		}
	}
	_, err := gc.client.PropagateGossip(ctx, &rpc.GossipMessage{
		SourceNode: message.SourceNode,
		StateMap:   stateMap,
	})
	return err
}

func (gc *GossipClient) PullState(ctx context.Context) (map[string]NodeState, error) {
	resp, err := gc.client.PullState(ctx, &rpc.PullRequest{})
	if err != nil {
		return nil, err
	}

	stateMap := make(map[string]NodeState)
	for k, v := range resp.StateMap {
		stateMap[k] = NodeState{
			NodeID:     v.NodeId,
			Status:     v.Status,
			Version:    v.Version,
			LastUpdate: parseTime(v.LastUpdate),
		}
	}
	return stateMap, nil
}

func parseTime(t string) time.Time {
	parsed, _ := time.Parse(time.RFC3339, t)
	return parsed
}
