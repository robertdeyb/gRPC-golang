package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"go-grpc/blog/blogpb"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "blog"
)

type server struct {
	conn *sql.DB
}

func (connection *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	db := connection.conn
	blog := req.GetBlog()
	authorId := blog.GetAuthorId()
	title := blog.GetTitle()
	content := blog.GetContent()
	sqlStatement := `INSERT INTO "blogs" (author_id, title, content) VALUES ($1, $2, $3)`
	if _, err := db.Exec(sqlStatement, authorId, title, content); err != nil {
		return nil, errors.Wrap(err, "Blog couldn't be inserted")
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       1,
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func main() {
	fmt.Println("Welcome to the server")
	lis, err := net.Listen("tcp", "0.0.0.0:50056")
	if err != nil {
		errors.Wrap(err, " Failed to listen the port")
	}
	s := grpc.NewServer()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		errors.Wrap(err, "Connection couldn't be opened")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		errors.Wrap(err, "Connection not established, ping didn't work")
	}
	blogpb.RegisterBlogServiceServer(s, &server{db})
	if err := s.Serve(lis); err != nil {
		errors.Wrap(err, "Failed to server the listener")
	}
}
