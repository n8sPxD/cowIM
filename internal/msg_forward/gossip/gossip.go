package gossip

import (
	"context"
	"github.com/n8sPxD/cowIM/internal/msg_forward/gossip/gossippb"
	"github.com/n8sPxD/cowIM/pkg/servicehub"
	"github.com/n8sPxD/cowIM/pkg/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func NewClient(neighbor string) *Client {
	conn, err := grpc.NewClient(neighbor, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// TODO: 重试，或者直接无视，从etcd中再次同步后再连接
		logx.Error("[UpdateNeighbors] Create client failed, error: ", err)
		return nil
	}
	client := gossippb.NewGossipClient(conn)
	return &Client{client, conn}
}

func (c *Client) Close() {
	c.conn.Close()
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

	Node   *Node
	discov *servicehub.DiscoveryHub
	regist *servicehub.RegisterHub

	id   int // worker id
	port int

	rounds int // 单次Gossip会传播给多少个节点
	depth  int // 单次Gossip会向下传染多少层
	retry  int // 传染失败后重试几次
}

// MustNewServer 创建 Gossip 服务。self 形式：<ip>:<port>
func MustNewServer(
	discov *servicehub.DiscoveryHub, regist *servicehub.RegisterHub,
	port, workerID, rounds, depth, retry int,
) *Server {
	ip, err := utils.GetLocalIP()
	if err != nil {
		panic(err)
	}
	s := &Server{
		UnimplementedGossipServer: &gossippb.UnimplementedGossipServer{},
		Node: &Node{
			self:      ip + ":" + strconv.Itoa(port),
			Data:      make(map[int32]Data),
			neighbors: make([]string, 0),
			clients:   make(map[string]*Client),
		},
		id:     workerID,
		port:   port,
		discov: discov,
		regist: regist,
		rounds: rounds,
		depth:  depth,
		retry:  retry,
	}
	return s
}

func (s *Server) Start() {
	ctx := context.Background()

	// 定时更新邻居
	go func() {
		s.UpdateNeighbors(ctx)

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.UpdateNeighbors(ctx)
			}
		}
	}()

	// 服务注册
	s.regist.Register(ctx, "gossip", s.port, uint16(s.id))

	// 启动 gRPC 服务
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))
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
	//logx.Debugf("[UpdateNeighbors] Updating neighbors from etcd, before neighbors: %v, self: %s", s.Node.neighbors, s.Node.self)
	var (
		neighbors = s.discov.GetServiceEndpoints(ctx, "gossip")
		exist     = make(map[string]bool)
	)

	s.Node.Lock()

	updated := make([]string, 0, len(neighbors))

	for _, neighbor := range neighbors {
		if neighbor == s.Node.self {
			continue
		} else {
			// 如果有新的节点加入，建立rpc客户端，加入映射表
			s.Node.clients[neighbor] = NewClient(neighbor)
		}
		exist[neighbor] = true
		updated = append(updated, neighbor)
	}

	for addr, neighbor := range s.Node.clients {
		if !exist[addr] {
			// 如果有节点退出，关闭rpc客户端，从映射表中删除
			neighbor.Close()
			delete(s.Node.clients, addr)
		}
	}

	s.Node.neighbors = updated

	s.Node.Unlock()

	//logx.Debugf("[UpdateNeighbors] Update success, current neighbors: %v, self: %s", s.Node.neighbors, s.Node.self)
}

// PushData 由邻居调用，邻居推送数据给自己，推模式更新
func (s *Server) PushData(ctx context.Context, in *gossippb.PushRequest) (*gossippb.PushResponse, error) {
	logx.Info("[PushData] Received data from neighbor, self: ", s.Node.self)
	s.Node.Lock()
	defer s.Node.Unlock()

	for _, remote := range in.Data {
		if _, exists := s.Node.Data[remote.Key]; !exists {
			s.Node.Data[remote.Key] = Data{remote.Value, remote.Version, remote.Timestamp}
		} else {
			local := s.Node.Data[remote.Key]
			if remote.Version > local.Version || (remote.Version == local.Version && remote.Timestamp > local.Timestamp) {
				s.Node.Data[remote.Key] = Data{remote.Value, remote.Version, remote.Timestamp}
			}
		}
	}

	s.Gossip(in)

	return &gossippb.PushResponse{}, nil
}

// RemoteUpdate 远程更新Data，自己同步
func (s *Server) RemoteUpdate(ctx context.Context, in *gossippb.RemoteRequest) (*gossippb.RemoteResponse, error) {
	logx.Debug("[RemoteUpdate] Received remote update, data: ", in.Data)
	s.Node.Lock()
	defer s.Node.Unlock()

	for _, remote := range in.Data {
		if _, exists := s.Node.Data[remote.Key]; !exists {
			s.Node.Data[remote.Key] = Data{remote.Value, remote.Version, remote.Timestamp}
		} else {
			local := s.Node.Data[remote.Key]
			if remote.Version > local.Version || (remote.Version == local.Version && remote.Timestamp > local.Timestamp) {
				s.Node.Data[remote.Key] = Data{remote.Value, remote.Version, remote.Timestamp}
			}
		}
	}

	go s.Gossip(&gossippb.PushRequest{Data: in.Data, Depth: 1})

	return &gossippb.RemoteResponse{}, nil
}

// Gossip 收到更新数据后，向邻居发送更新
func (s *Server) Gossip(updates *gossippb.PushRequest) {
	logx.Debugf("[Gossip] Received gossip request, self: %s", s.Node.self)
	if updates.Depth >= int32(s.depth) {
		return
	}

	s.Node.RLock()
	defer s.Node.RUnlock()

	// 可以加延迟，这里设定收到消息后立刻发送消息
	// 加延迟后，本地需要维护 chan 缓存信息，
	// 然后定时从 chan中取出所有消息再进行Gossip

	// 没有邻居，不需要更新
	if len(s.Node.neighbors) <= 0 {
		logx.Debug("[Gossip] No neighbors found! self: ", s.Node.self)
		return
	}

	for range s.rounds {
		var (
			luckyboy = rand.Intn(len(s.Node.neighbors))
			current  = s.Node.neighbors[luckyboy]
			request  = gossippb.PushRequest{Data: updates.Data, Depth: updates.Depth + 1}
		)
		var pushUpdate func(int)
		pushUpdate = func(retry int) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			if retry == 0 {
				return
			} else if _, err := s.Node.clients[current].PushData(ctx, &request); err != nil {
				logx.Error("[Gossip] Push data to neighbor failed, error: ", err)
				logx.Infof("[Gossip] Retry push data to neighbor, retry: %d", retry)
				time.Sleep(1 * time.Second)
				pushUpdate(retry - 1)
			}
		}
		pushUpdate(s.retry)
	}
}
