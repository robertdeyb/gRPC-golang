package main

import (
	"context"
	"fmt"
	"go-grpc/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := blogpb.NewBlogServiceClient(conn)

	//CREATE USER PROFILE

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       1,
			AuthorId: "Dave",
			Title:    "Sample Title",
			Content:  "First Content",
		},
	}
	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while creating blog %v", err)
	}
	fmt.Println(res)
}
