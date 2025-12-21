package main

import (
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	pb "github.com/lavish440/go-microservice/calculator/proto"
)

var serverName = os.Getenv("SERVER_NAME")

type server struct {
	pb.UnimplementedCalcServiceServer
}

func (s *server) Add(ctx context.Context, req *pb.CalcRequest) (*pb.CalcResponse, error) {
	return &pb.CalcResponse{Result: req.A + req.B, ServerName: serverName}, nil
}

func (s *server) Sub(ctx context.Context, req *pb.CalcRequest) (*pb.CalcResponse, error) {
	return &pb.CalcResponse{Result: req.A - req.B, ServerName: serverName}, nil
}

func (s *server) Mul(ctx context.Context, req *pb.CalcRequest) (*pb.CalcResponse, error) {
	return &pb.CalcResponse{Result: req.A * req.B, ServerName: serverName}, nil
}

func (s *server) Div(ctx context.Context, req *pb.CalcRequest) (*pb.CalcResponse, error) {
	if req.B == 0 {
		return nil, status.Error(codes.InvalidArgument, "division by zero is not allowed")
	}

	return &pb.CalcResponse{Result: req.A / req.B, ServerName: serverName}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterCalcServiceServer(grpcServer, &server{})

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(
		"calculator.CalcService",
		grpc_health_v1.HealthCheckResponse_SERVING,
	)

	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	log.Println("Calculator gRPC server running on :50051")
	log.Fatal(grpcServer.Serve(lis))
}
