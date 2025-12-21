package main

import (
	"net"
	"time"
)

func checkHealth(b *Backend) {
	conn, err := net.DialTimeout("tcp", b.Addr, 2*time.Second)
	if err != nil {
		b.Healthy = false
		return
	}
	defer conn.Close()

	b.Healthy = true
}
