package main

import (
	"log"
)

type RoundRobin struct {
	backends []Backend
	weights  []int
	last     int
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (rr *RoundRobin) Update(backends []Backend) {
	rr.backends = backends
	rr.weights = []int{}

	for i := range backends {
		for w := 0; w < backends[i].Weight; w++ {
			rr.weights = append(rr.weights, i)
		}
	}
	log.Printf("[lb] Updated weighted Round Robin with %d entries: %v\n", len(rr.weights), rr.weights)
}

func (rr *RoundRobin) Next() *Backend {
	if len(rr.weights) == 0 || len(rr.backends) == 0 {
		return nil
	}

	rr.last = (rr.last + 1) % len(rr.weights)
	selectedIdx := rr.weights[rr.last]
	return &rr.backends[selectedIdx]
}
