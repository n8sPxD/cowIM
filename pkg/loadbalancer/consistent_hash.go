package loadbalancer

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// ConsistentHash 一致性哈希法实现
type ConsistentHash struct {
	hashRing    []uint32
	nodes       map[uint32]string
	virtualNode int
	mu          sync.RWMutex
}

// NewConsistentHash 创建一致性哈希实例
func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{
		nodes:       make(map[uint32]string),
		virtualNode: 3, // 默认每个节点对应3个虚拟节点
	}
}

// hash 使用crc32生成哈希值
func (ch *ConsistentHash) hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// addNode 添加节点到哈希环
func (ch *ConsistentHash) addNode(node string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	// 为每个物理节点创建虚拟节点
	for i := 0; i < ch.virtualNode; i++ {
		virtualKey := node + "#" + strconv.Itoa(i)
		hash := ch.hash(virtualKey)
		ch.hashRing = append(ch.hashRing, hash)
		ch.nodes[hash] = node
	}
	sort.Slice(ch.hashRing, func(i, j int) bool {
		return ch.hashRing[i] < ch.hashRing[j]
	})
}

// removeNode 从哈希环中移除节点
func (ch *ConsistentHash) removeNode(node string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	// 移除对应的虚拟节点
	for i := 0; i < ch.virtualNode; i++ {
		virtualKey := node + "#" + strconv.Itoa(i)
		hash := ch.hash(virtualKey)
		delete(ch.nodes, hash)

		// 移除哈希环中的值
		for idx, val := range ch.hashRing {
			if val == hash {
				ch.hashRing = append(ch.hashRing[:idx], ch.hashRing[idx+1:]...)
				break
			}
		}
	}
}

// Take 根据一致性哈希分配一个节点
func (ch *ConsistentHash) Take(endpoints []string) string {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	// 如果当前哈希环为空，重新初始化
	if len(ch.hashRing) == 0 {
		for _, endpoint := range endpoints {
			ch.addNode(endpoint)
		}
	}

	// 如果仍然没有可用节点，返回空字符串
	if len(ch.hashRing) == 0 {
		return ""
	}

	// 随机选择一个哈希值（例如用当前时间戳或其他请求标识符）
	// 示例中用简单字符串"key"进行哈希，可以用请求IP或其他标识代替
	requestKey := "key" // 替换为具体请求标识符
	hash := ch.hash(requestKey)

	// 查找大于等于此哈希值的节点（顺时针第一个节点）
	idx := sort.Search(len(ch.hashRing), func(i int) bool {
		return ch.hashRing[i] >= hash
	})

	// 如果超出范围，回到哈希环的第一个节点
	if idx == len(ch.hashRing) {
		idx = 0
	}

	return ch.nodes[ch.hashRing[idx]]
}
