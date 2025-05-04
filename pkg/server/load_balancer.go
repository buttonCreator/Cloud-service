package server

import (
	"errors"
	"sync/atomic"
)

type BalancerStrategy int

const (
	RoundRobin BalancerStrategy = iota
	LeastConnections
)

var strategyCreators = map[BalancerStrategy]func() BalanceStrategy{
	RoundRobin:       func() BalanceStrategy { return &RoundRobinStrategy{} },
	LeastConnections: func() BalanceStrategy { return &LeastConnectionsStrategy{} },
}

type RoundRobinStrategy struct {
	counter uint64
}

type LeastConnectionsStrategy struct {
	counters []uint64
}

func (r *RoundRobinStrategy) NextBackend(backends []*Backend) (*Backend, error) {
	next := atomic.AddUint64(&r.counter, 1)
	length := uint64(len(backends))
	index := next % length

	for i := uint64(0); i < length; i++ {
		idx := (index + i) % length
		if backends[idx].Alive {
			return backends[idx], nil
		}
	}

	return nil, errors.New("no available backends")
}

func (lc *LeastConnectionsStrategy) NextBackend(backends []*Backend) (*Backend, error) {
	var best *Backend
	minValue := ^uint64(0)

	for i, b := range backends {
		if !b.Alive {
			continue
		}

		current := atomic.LoadUint64(&lc.counters[i])
		if current < minValue {
			minValue = current
			best = b
		}
	}

	if best == nil {
		return nil, errors.New("no available backends")
	}

	atomic.AddUint64(&lc.counters[len(lc.counters)-1], 1)
	return best, nil
}

// UnmarshalText implements TextUnmarshaler
func (e *BalancerStrategy) UnmarshalText(text []byte) error {
	switch string(text) {
	case "round_robin":
		*e = RoundRobin
	case "least_connections":
		*e = LeastConnections
	}

	return nil
}

func (b *BalancerStrategy) String() string {
	if *b == RoundRobin {
		return "round_robin"
	}

	return "least_connections"
}
