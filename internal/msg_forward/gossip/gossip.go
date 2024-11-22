package gossip

import (
	"context"
	"github.com/n8sPxD/cowIM/internal/msg_forward/gossip/gossippb"
	"github.com/n8sPxD/cowIM/pkg/servicehub"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

type Data struct {
	Value     int32
	Version   int64
	Timestamp int64
}

type Client struct {
	gossippb.GossipClient
	conn *grpc.ClientConn
}

type Node struct {
	sync.RWMutex
	self      string
	Data      map[int32]Data
	neighbors []string
	clients   map[string]*Client
}

type Server struct {
	*gossippb.UnimplementedGossipServer
	node   *Node
	discov *servicehub.DiscoveryHub

	rounds int // 单次Gossip会传播给多少个节点
	depth  int // 单次Gossip会向下传染多少层
	retry  int // 传染失败后重试几次
}

// MustNewServer 创建 Gossip 服务。self 形式：<ip>:<port>
func MustNewServer(hub *servicehub.DiscoveryHub, self string, rounds, depth, retry int) *Server {
	s := &Server{
		node: &Node{
			self:      self,
			Data:      make(map[int32]Data),
			neighbors: make([]string, 0),
			clients:   make(map[string]*Client),
		},
		discov: hub,
		rounds: rounds,
		depth:  depth,
		retry:  retry,
	}
	return s
}

func (s *Server) Start(ctx context.Context, port int) {
	s.UpdateNeighbors(ctx)

	// 定时更新邻居
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		select {
		case <-ticker.C:
			s.UpdateNeighbors(ctx)
		}
	}()

	// 启动 gRPC 服务
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		logx.Error("[Start] Listen failed, error: ", err)
		return
	}
	rpcserver := grpc.NewServer()
	gossippb.RegisterGossipServer(rpcserver, s)
	if err := rpcserver.Serve(lis); err != nil {
		logx.Error("[Start] Serve failed, error: ", err)
		return
	}
}

// UpdateNeighbors 更新邻居
func (s *Server) UpdateNeighbors(ctx context.Context) {
	var (
		neighbors = s.discov.GetServiceEndpoints(ctx, "msgfwd")
		exist     = make(map[string]bool)
	)

	s.node.Lock()
	defer s.node.Unlock()

	s.node.neighbors = neighbors

	for _, neighbor := range neighbors {
		if neighbor == s.node.self {
			continue
		} else {
			// 如果有新的节点加入，建立rpc客户端，加入映射表
			if s.node.clients[neighbor] == nil {
				conn, err := grpc.NewClient(neighbor)
				if err != nil {
					// TODO: 重试，或者直接无视，从etcd中再次同步
					logx.Error("[UpdateNeighbors] Create client failed, error: ", err)
					continue
				}
				client := gossippb.NewGossipClient(conn)
				s.node.clients[neighbor] = &Client{client, conn}
			}
			exist[neighbor] = true
		}
	}

	for addr, neighbor := range s.node.clients {
		if !exist[addr] {
			// 如果有节点退出，关闭rpc客户端，从映射表中删除
			neighbor.conn.Close()
			delete(s.node.clients, addr)
		}
	}
}

// PushData 由邻居调用，邻居推送数据给自己，推模式更新
func (s *Server) PushData(ctx context.Context, in *gossippb.PushRequest) (*gossippb.PushResponse, error) {
	s.node.Lock()
	defer s.node.Unlock()

	for _, remote := range in.Data {
		if _, exists := s.node.Data[remote.Key]; !exists {
			s.node.Data[remote.Key] = Data{remote.Value, remote.Version, remote.Timestamp}
		} else {
			local := s.node.Data[remote.Key]
			if remote.Version > local.Version || (remote.Version == local.Version && remote.Timestamp > local.Timestamp) {
				s.node.Data[remote.Key] = Data{remote.Value, remote.Version, remote.Timestamp}
			}
		}
	}

	s.Gossip(ctx, in)

	return &gossippb.PushResponse{}, nil
}

// RemoteUpdate 远程更新Data，自己同步
func (s *Server) RemoteUpdate(ctx context.Context, in *gossippb.RemoteRequest) (*gossippb.RemoteResponse, error) {
	s.node.Lock()
	defer s.node.Unlock()

	for _, remote := range in.Data {
		if _, exists := s.node.Data[remote.Key]; !exists {
			s.node.Data[remote.Key] = Data{remote.Value, remote.Version, remote.Timestamp}
		} else {
			local := s.node.Data[remote.Key]
			if remote.Version > local.Version || (remote.Version == local.Version && remote.Timestamp > local.Timestamp) {
				s.node.Data[remote.Key] = Data{remote.Value, remote.Version, remote.Timestamp}
			}
		}
	}

	go s.Gossip(ctx, &gossippb.PushRequest{Data: in.Data, Depth: 1})

	return &gossippb.RemoteResponse{}, nil
}

// Gossip 收到更新数据后，向邻居发送更新
func (s *Server) Gossip(ctx context.Context, updates *gossippb.PushRequest) {
	if updates.Depth >= int32(s.depth) {
		return
	}

	s.node.RLock()
	defer s.node.RUnlock()

	// 可以加延迟，这里设定收到消息后立刻发送消息
	// 加延迟后，本地需要维护 chan 缓存信息，
	// 然后定时从 chan中取出所有消息再进行Gossip

	for range s.rounds {
		var (
			luckyboy = rand.Intn(len(s.node.neighbors))
			current  = s.node.neighbors[luckyboy]
			request  = gossippb.PushRequest{Data: updates.Data, Depth: updates.Depth + 1}
		)
		var pushUpdate func(int)
		pushUpdate = func(retry int) {
			if retry == 0 {
				return
			} else if _, err := s.node.clients[current].PushData(ctx, &request); err != nil {
				logx.Error("[Gossip] Push data to neighbor failed, error: ", err)
				logx.Infof("[Gossip] Retry push data to neighbor, retry: %d", retry)
				pushUpdate(retry - 1)
			}
		}
		go pushUpdate(s.retry)
	}
}
