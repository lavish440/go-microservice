package main

import (
	"io"
	"log"
	"net"
	"time"
)

func startTCPLB(rr *RoundRobin) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("TCP LB listening on :50051")

	for {
		clientConn, err := lis.Accept()
		if err != nil {
			log.Printf("Failed to accept client connection: %v", err)
			continue
		}

		activeConnections.Inc()
		go func() {
			defer activeConnections.Dec()
			start := time.Now()

			backend := rr.Next()
			if backend == nil || !backend.Healthy {
				log.Println("No healthy backend available")
				clientConn.Close()
				return
			}

			totalRequests.WithLabelValues(backend.Addr).Inc()
			log.Printf("Forwarding to backend: %s", backend.Addr)

			backendConn, err := net.Dial("tcp", backend.Addr)
			if err != nil {
				log.Printf("Failed to connect to backend: %v", err)
				clientConn.Close()
				return
			}

			connectionLatency.Observe(time.Since(start).Seconds())
			go func() {
				defer clientConn.Close()
				defer backendConn.Close()
				io.Copy(backendConn, clientConn)
			}()
			go func() {
				defer clientConn.Close()
				defer backendConn.Close()
				io.Copy(clientConn, backendConn)
			}()
		}()
	}
}
