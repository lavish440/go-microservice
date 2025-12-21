package main

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func checkHealth(b *Backend) {
	conn, err := grpc.NewClient(b.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Healthy = false
		return
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "calculator.CalcService"})
	if err != nil || resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		b.Healthy = false
		return
	}

	b.Healthy = true
}
