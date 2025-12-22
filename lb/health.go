package main

import (
	"log"
	"net"
	"time"
)

func checkHealth(backend *Backend) {
	for {
		conn, err := net.DialTimeout("tcp", backend.Addr, 2*time.Second)
		if err != nil {
			backend.Healthy = false
			log.Printf("[health] Backend %s marked as unhealthy: %v", backend.Addr, err)
		} else {
			backend.Healthy = true
			conn.Close()
		}

		time.Sleep(5 * time.Second)
	}
}
