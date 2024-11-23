package loadbalancer

import "sync/atomic"

// RoundRobin 轮询法
type RoundRobin struct {
	acc int64
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (rr *RoundRobin) Take(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	n := atomic.AddInt64(&rr.acc, 1)
	idx := int(n % int64(len(endpoints)))
	return endpoints[idx]
}
