package main

import (
	"context"
	pb "github.com/shaojunda/Go-000/Week4/explorer/api/book/v1"
	"google.golang.org/grpc"
	"log"
	"time"

	//pb "github.com/shaojunda/Go-000/Week4/explorer/api/book/v1"
)

const (
	address = "localhost:8787"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBookServiceClient(conn)
	searchBook(c)
	time.Sleep(3 * time.Second)
	createBook(c)
}

func searchBook(c pb.BookServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	res, err := c.SearchBook(ctx, &pb.SearchBookRequest{
		QueryString: "haha",
	})
	if err != nil {
		log.Fatal("could not fund")
	}
	log.Printf("gRPC response: %s\n", res.String())
}

func createBook(c pb.BookServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	res, err := c.CreateBook(ctx, &pb.CreateBookRequest{
		Name:      "",
		Author:    "",
		Publisher: "",
		Price:     "1000",
		Url:       "",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("gRPC response: %s\n", res.String())
}