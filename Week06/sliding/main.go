package main

import (
	"context"
	"errors"
	"github.com/shaojunda/Go-000/Week06/sliding/bucket"
	"golang.org/x/sync/errgroup"
	"log"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, gCtx := errgroup.WithContext(ctx)
	rw := bucket.NewRollingTimeWindow(100)
	g.Go(func() error {
		return rw.Check(gCtx)
	})

	g.Go(func() error {
		for {
			if !rw.IsLimit {
				rw.Mutex.Lock()
				rw.Counter++
				rw.Mutex.Unlock()
				log.Println("processed")
				rand.Seed(time.Now().UnixNano())
				n := rand.Intn(10)
				time.Sleep(time.Duration(n) * time.Millisecond)
			} else {
				return errors.New("limiting")
			}
		}
	})


	err := g.Wait()
	if err == nil {
		log.Println(err)
	}
}