package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"go-grpc/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type config struct {
	db struct {
		dsn string
	}
}

type server struct{}

type Blog struct {
	ID       int    `json:"id"`
	AuthorID string `json: author_id`
	Content  string `json: content`
	Title    string `json: title`
}

func main() {
	//if we crash something, we get the file and error number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var cfg config
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://postgres:root@localhost/blogs?sslmode=disable", "Postgres connection string")
	db, err := openDB(cfg)
	fmt.Println("Connecting to Posgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Blog Service Started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server..")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until the signal is received
	<-ch
	fmt.Println("Stoppping the server")
	s.Stop()
	fmt.Println("Stopping the listener")
	lis.Close()
	fmt.Println("Closing Posgres DB")
	db.Close()
	fmt.Println("End of Program")

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}
