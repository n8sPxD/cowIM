package loadbalancer

const (
	RoundRobinBalancer = iota
	ConsistentHashBalancer
)

type LoadBalancer interface {
	Take([]string) string
}

func NewLoadBalancer(lbtype int) LoadBalancer {
	switch lbtype {
	case RoundRobinBalancer:
		return &RoundRobin{}
	case ConsistentHashBalancer:
		return &ConsistentHash{}
	default:
		return nil
	}
}
