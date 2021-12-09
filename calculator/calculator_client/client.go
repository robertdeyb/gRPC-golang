package main

import (
	"context"
	"fmt"
	"go-grpc/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	doErrorUnary(c, 10)
	doErrorUnary(c, -2)

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

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecompositionServer Streaming RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happend %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Compute Average Client Streaming RPC...")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}
	numbers := []int32{3, 5, 9, 54, 23}
	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}
	fmt.Printf("The average is: %v", res.GetAverage())

}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Find Maximum Bi Directional Streaming RPC...")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}
	waitc := make(chan struct{})

	// send a bunch of message to the client (go routine)
	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Println("Sending message: %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive a bunch of messages from the client (go routine)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received a new maximum: %v\n", res.GetMaximumNumber())

		}
		close(waitc)
	}()
	//block until everything is done
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient, n int32) {
	fmt.Println("Starting to do a SquareRoot Unary RPC..")
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: n,
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from grpc
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
		}
	}
	fmt.Println("Result of square root of %v: %v", n, res.GetNumberRoot())
}
