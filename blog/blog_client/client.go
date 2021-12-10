package main

import (
	"context"
	"fmt"
	"go-grpc/blog/blogpb"
	"io"
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

	fmt.Println("Create blog request")
	doCreate(c)

	//GET Blog
	fmt.Println("Read blog request")
	doRead(c)

	update Blog
	fmt.Println("Update blog request")
	doUpdate(c)

	delete Blog
	fmt.Println("Delete blog request")
	doDelete(c, 7)

	list Blogs
	fmt.Println("List blog request")
	doList(c)

}

func doCreate(c blogpb.BlogServiceClient) {
	//CREATE Blog
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Dave Torrente",
			Title:    "Sample Title",
			Content:  "Sample Content",
		},
	}
	createdBlogRes, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while creating blog %v", err)
	}
	fmt.Println("Blog has been created: %v", createdBlogRes)
}

func doRead(c blogpb.BlogServiceClient) {
	// //test id
	id := int64(2)
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: id,
	})
	if err2 != nil {
		fmt.Sprintf("Error happened while reading: %v", err2)
		return
	}
	readBlogReq := &blogpb.ReadBlogRequest{BlogId: id}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readBlogErr)
	}
	fmt.Printf("Blog was read: %v \n", readBlogRes)
}

func doUpdate(c blogpb.BlogServiceClient) {
	newBlog := &blogpb.Blog{
		Id:       5,
		AuthorId: "Changed Author Again",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)
	fmt.Println(updateRes)
}
func doDelete(c blogpb.BlogServiceClient, id int64) {
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: id})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)
}
func doList(c blogpb.BlogServiceClient) {
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
