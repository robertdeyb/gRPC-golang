package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"google.golang.org/grpc"

	"go-grpc/blog/blogpb"
)

func TestCreateUserProfile(t *testing.T) {
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := blogpb.NewBlogServiceClient(conn)
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       2,
			AuthorId: "Robert",
			Title:    "Sample Title Test",
			Content:  "Sample Test",
		},
	}
	res, err := c.CreateBlog(context.Background(), req)
	fmt.Println(res)
}
