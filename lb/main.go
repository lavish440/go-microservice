package main

import (
	"log"
	"os"
	"time"
)

func main() {
	rr := NewRoundRobin()

	service := os.Getenv("GRPC_SERVICE")
	if service == "" {
		log.Fatal("Environment variable GRPC_SERVICE not set!")
	}

	go func() {
		for {
			backends := discoverFromDocker(service)

			rr.Update(backends)

			for _, key := range rr.keys {
				go checkHealth(rr.backends[key])
			}

			time.Sleep(10 * time.Second)
		}
	}()

	startTCPLB(rr)
}
