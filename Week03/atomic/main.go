package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type Config struct {
	a []int
}

func (c *Config) T() {}

func main() {
	var v atomic.Value
	v.Store(&Config{})

	go func() {
		i := 0
		for {
			i++
			cfg := &Config{a: []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			v.Store(cfg)
		}
	}()

	ctx := context.Background()
	ctx.Done()

	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < 10; n++ {
				cfg := v.Load().(*Config)
				fmt.Printf("%v\n", cfg)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
