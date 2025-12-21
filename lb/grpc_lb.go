package main

import (
	"context"
	"log"
	"net"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func startGRPCLB(rr *RoundRobin) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("gRPC LB listening on :50051")

	grpcServer := grpc.NewServer()

	director := func(ctx context.Context, _ string) (context.Context, grpc.ClientConnInterface, error) {
		backend := rr.Next()
		if backend == nil {
			return ctx, nil, status.Error(codes.Unavailable, "no healthy backends")
		}

		conn, err := grpc.NewClient(backend.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("Failed to dial backend:", err)
			return ctx, nil, err
		}

		return ctx, conn, nil
	}

	proxy.RegisterService(grpcServer, director, "calculator.CalcService", "Add", "Sub", "Mul", "Div")

	log.Fatal(grpcServer.Serve(lis))
}
