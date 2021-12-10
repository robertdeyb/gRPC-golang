package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"go-grpc/blog/blogpb"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type blogItem struct {
	ID       int    `json:"id"`
	AuthorID string `json:"author_id"`
	Content  string `json:"content"`
	Title    string `json:"title"`
}

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
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	query := `INSERT INTO "blogs" (author_id, title, content) VALUES ($1, $2, $3) RETURNING id`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	lastInsertId := int64(0)
	if err = stmt.QueryRow(data.AuthorID, data.Title, data.Content).Scan(&lastInsertId); err != nil {
		return nil, errors.Wrap(err, "Blog couldn't be inserted")
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       lastInsertId,
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (connection *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {

	db := connection.conn
	id := req.GetBlogId()
	sqlStatement := `select author_id, title, content from "blogs" where id=$1`
	var author_id, title, content string
	err := db.QueryRow(sqlStatement, id).Scan(&author_id, &title, &content)
	if err != nil {
		errors.Wrap(err, "Blog couldn't be returned")
	}
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       id,
			AuthorId: author_id,
			Title:    title,
			Content:  content,
		},
	}, nil
}

func (connection *server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("Update blog request")
	blog := req.GetBlog()

	// create an empty struct
	data := &blogItem{}

	// we update our internal struct
	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()
	data.ID = int(blog.GetId())
	db := connection.conn
	sqlStatement := `UPDATE "blogs" SET author_id=$1, title=$2,content=$3 WHERE "id" =$4;`
	if _, err := db.Exec(sqlStatement, data.AuthorID, data.Content, data.Title, data.ID); err != nil {
		return nil, err
	}
	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       int64(data.ID),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil

}

func (connection *server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("Delete blog request")
	db := connection.conn
	id := req.GetBlogId()
	sqlStatement := `delete from "blogs" where id=$1`
	if _, err := db.Exec(sqlStatement, id); err != nil {
		errors.Wrap(err, "Blog couldn't be deleted")
	}
	return &blogpb.DeleteBlogResponse{BlogId: req.GetBlogId()}, nil
}

func (connection *server) ListBlog(_ *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("List blog request")
	db := connection.conn
	sqlStatement := `select * from "blogs"`
	result, err := db.Query(sqlStatement)
	defer result.Close()
	if err != nil {
		fmt.Println(err)
	}
	for result.Next() {
		var id int64
		var author_id, title, content string
		if err = result.Scan(&author_id, &title, &content, &id); err != nil {
			errors.Wrap(err, "Blogs couln't be listed")
		}

		stream.Send(&blogpb.ListBlogResponse{Blog: &blogpb.Blog{
			AuthorId: author_id,
			Title:    title,
			Content:  content,
			Id:       id,
		}})
	}
	return nil
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
