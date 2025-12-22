package main

import (
	"io"
	"log"
	"os"
	"time"
)

func main() {
	rr := NewRoundRobin()

	service := os.Getenv("GRPC_SERVICE")
	if service == "" {
		log.Fatal("GRPC_SERVICE Environment variable not set!")
	}

	if os.Getenv("DEBUG") != "true" {
		log.SetOutput(io.Discard)
	}

	go func() {
		backends := []Backend{}
		for {
			backends = discoverFromDocker(service, backends)
			rr.Update(backends)
			for i := range backends {
				go checkHealth(&backends[i])
			}
			time.Sleep(10 * time.Second)
		}
	}()

	go startMetricsEndpoint()

	startTCPLB(rr)
}
