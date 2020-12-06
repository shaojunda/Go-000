package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	tr := NewTracer()
	go tr.Run()
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(20*time.Second))
	defer cancel()
	tr.Shutdown(ctx)
}

func NewTracer() *Tracer {
	return &Tracer{
		ch: make(chan string, 10),
	}
}

type Tracer struct {
	ch   chan string
	stop chan struct{}
}

func (t *Tracer) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracer) Run() {
	for data := range t.ch {
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}
	t.stop <- struct{}{}
}

func (t *Tracer) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
		fmt.Println("chan stopped")
	case <-ctx.Done():
		fmt.Println("context done")
	}
}
