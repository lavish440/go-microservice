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

			handleTCPConnection(clientConn, rr)

			connectionLatency.Observe(time.Since(start).Seconds())
		}()
	}
}

func handleTCPConnection(clientConn net.Conn, rr *RoundRobin) {
	backend := rr.Next()
	if backend == nil {
		log.Println("No healthy backend available")
		clientConn.Close()
		return
	}

	totalRequests.WithLabelValues(backend.Addr).Inc()

	// log.Printf("Forwarding TCP connection to backend: %s", backend.Addr)

	backendConn, err := net.Dial("tcp", backend.Addr)
	if err != nil {
		log.Printf("Failed to connect to backend %s: %v", backend.Addr, err)
		clientConn.Close()
		return
	}

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
}
