package main

import (
	"context"
	"fmt"
	"go-grpc/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculator(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Println("Calculator function was invoked with %v\n", req)
	firstNum := req.GetCalculate().FirstNum
	lastNum := req.GetCalculate().LastNum
	//response
	result := firstNum + lastNum
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}
	return res, nil
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
