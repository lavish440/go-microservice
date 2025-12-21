package main

import (
	"os"
)

func main() {
	rr := NewRoundRobin()

	go func() {
		for {
			backends := discoverFromDocker(os.Getenv("GRPC_SERVICE"))
			rr.Update(backends)
			for _, key := range rr.keys {
				go checkHealth(rr.backends[key])
			}
		}
	}()

	startGRPCLB(rr)
}
