package lb

import (
	"sync/atomic"
)

const (
	ROUNDROBIN = iota
)

type LoadBalancer interface {
	Take([]string) string
}

func NewLoadBalancer(lbtype int) LoadBalancer {
	switch lbtype {
	case ROUNDROBIN:
		return &RoundRobin{}
	default:
		return nil
	}
}

// RoundRobin 轮询法
type RoundRobin struct {
	acc int64
}

func (rr *RoundRobin) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	n := atomic.AddInt64(&rr.acc, 1)
	idx := int(n % int64(len(endpoints)))
	return endpoints[idx]
}
