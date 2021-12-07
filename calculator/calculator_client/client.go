package main

import (
	"context"
	"fmt"
	"go-grpc/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	doUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.CalculatorRequest{
		Calculate: &calculatorpb.Calculate{
			FirstNum: 3,
			LastNum:  10,
		},
	}
	res, err := c.Calculator(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculator RPC: %v", err)
	}
	log.Printf("Response from Calculator. %v", res.Result)
}
