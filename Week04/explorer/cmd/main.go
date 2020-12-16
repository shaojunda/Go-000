package main

import (
	"context"
	"github.com/pkg/errors"
	pb "github.com/shaojunda/Go-000/Week4/explorer/api/book/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

const (
	port = ":8787"
)

func main() {
	var ctx context.Context
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	defer cancel()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	g.Go(func() error {
		bookService := InitializeBookService()
		lis, err := net.Listen("tcp", port)
		if err != nil {
			return errors.WithMessage(err, "failed to listen")
		}
		s := grpc.NewServer()
		log.Println("gRPC server is running.")

		pb.RegisterBookServiceServer(s, bookService)
		if err := s.Serve(lis); err != nil {
			return errors.WithMessage(err, "failed to serve: %v")
		}
		return nil
	})

	g.Go(func() error {
		select {
		case <-c:
			log.Printf("cancel sig received.")
			cancel()
			return errors.New("cancel sig received")
		case <-ctx.Done():
			log.Printf("context done.")
			return errors.New("cancel sig received")
		}
	})

	if err := g.Wait(); err == nil {
		log.Println("start server successfully")
	}
}
