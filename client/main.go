package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/lavish440/go-microservice/calculator/proto"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewCalcServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("gRPC Calculator Client")
	fmt.Println("Operations: add, sub, mul, div")
	fmt.Println("Type 'exit' to quit")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("bye")
			return
		}

		parts := strings.Fields(input)
		if len(parts) != 3 {
			fmt.Println("Usage: <op> <a> <b>")
			continue
		}

		op := parts[0]
		a, err1 := strconv.ParseFloat(parts[1], 64)
		b, err2 := strconv.ParseFloat(parts[2], 64)

		if err1 != nil || err2 != nil {
			fmt.Println("Invalid numbers")
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var res *pb.CalcResponse
		var err error

		switch op {
		case "add":
			res, err = client.Add(ctx, &pb.CalcRequest{A: a, B: b})
		case "sub":
			res, err = client.Sub(ctx, &pb.CalcRequest{A: a, B: b})
		case "mul":
			res, err = client.Mul(ctx, &pb.CalcRequest{A: a, B: b})
		case "div":
			res, err = client.Div(ctx, &pb.CalcRequest{A: a, B: b})
		default:
			fmt.Println("Unknown operation")
			continue
		}

		if err != nil {
			if st, ok := status.FromError(err); ok {
				fmt.Printf("Error: %s (%s)\n", st.Message(), st.Code())
			} else {
				fmt.Println("Error:", err)
			}
			continue
		}

		fmt.Println(res.ServerName+":", res.Result)
	}
}
