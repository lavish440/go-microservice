package main

import "sync"

type RoundRobin struct {
	mu       sync.Mutex
	backends map[string]*Backend
	keys     []string
	index    int
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		backends: make(map[string]*Backend),
		keys:     []string{},
		index:    0,
	}
}

func (rr *RoundRobin) Update(newBackends []Backend) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	newKeys := make([]string, 0, len(newBackends))
	for _, b := range newBackends {
		if existing, ok := rr.backends[b.Addr]; ok {
			existing.Addr = b.Addr
		} else {
			rr.backends[b.Addr] = &Backend{Addr: b.Addr, Healthy: false}
		}
		newKeys = append(newKeys, b.Addr)
	}
	rr.keys = newKeys
}

func (rr *RoundRobin) Next() *Backend {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	n := len(rr.keys)
	if n == 0 {
		return nil
	}

	for range n {
		key := rr.keys[rr.index]
		rr.index = (rr.index + 1) % n
		if backend := rr.backends[key]; backend != nil && backend.Healthy {
			return backend
		}
	}
	return nil
}
