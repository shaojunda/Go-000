package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, gCtx := errgroup.WithContext(ctx)

	// random cancel
	randomCancel(ctx, cancel)

	// start walking server
	g.Go(func() error {
		return startWalkingServer(gCtx)
	})

	// start eating server
	g.Go(func() error {
		return startEatingServer(gCtx)
	})

	// start sleeping server
	g.Go(func() error {
		return startSleepingServer(gCtx)
	})

	// start bang server
	g.Go(func() error {
		return startBangServer(gCtx)
	})

	g.Go(func() error {
		return listenSig(ctx, cancel)
	})

	err := g.Wait()

	if err == nil {
		fmt.Println("successfully exit")
	} else {
		fmt.Printf("group error %+v", err)
	}
	for {
		select {
		case <-ctx.Done():
			countdown(context.Background(), 10)
			log.Println("All services are closed successfully.")
			return
		}
	}
}

func countdown(ctx context.Context, seconds int) {
	tick := time.Tick(1 * time.Second)
	for countdown := seconds; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-ctx.Done():
			return
		case <-tick:
		}
	}
}

func startBangServer(ctx context.Context) error {
	m := &http.ServeMux{}
	s := http.Server{Addr: ":8093", Handler: m}
	log.Println("bang server listen on 8093")

	m.HandleFunc("/bang", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "bang!\n")
		if err != nil {
			log.Fatalf("bang error: %+v", err)
		}
		timeout, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		s.Shutdown(timeout)
	})
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("bang server panic %+")
			}
		}()
		<-ctx.Done()
		log.Println("close signal received by bang server.")
		timeout, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		s.Shutdown(timeout)
	}()

	return s.ListenAndServe()
}

func startSleepingServer(ctx context.Context) error {
	m := &http.ServeMux{}
	s := http.Server{Addr: ":8092", Handler: m}
	m.HandleFunc("/sleeping", sleepHandler)
	log.Println("sleeping server listen on 8092")

	go func() {
		<-ctx.Done()
		log.Println("close signal received by sleeping server.")
		timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer func() {
			cancel()
			if err := recover(); err != nil {
				log.Printf("sleeping server panic %+v\n", err)
			}
		}()

		s.Shutdown(timeout)
		log.Println("sleeping server closed.")
	}()

	return s.ListenAndServe()
}

func startEatingServer(ctx context.Context) error {
	m := &http.ServeMux{}
	s := http.Server{Addr: ":8091", Handler: m}
	m.HandleFunc("/eating", eatHandler)
	log.Println("eat server listen on 8091")

	go func() {
		<-ctx.Done()
		log.Println("close signal received by eating server.")
		timeout, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer func() {
			cancel()
			if err := recover(); err != nil {
				log.Printf("eating server panic %+v\n", err)
			}
		}()

		s.Shutdown(timeout)
		log.Println("eating server closed.")
	}()

	return s.ListenAndServe()
}

func startWalkingServer(ctx context.Context) error {
	m := &http.ServeMux{}
	s := http.Server{Addr: ":8090", Handler: m}
	m.HandleFunc("/walking", walkHandler)
	log.Println("walk server listen on 8090")

	go func() {
		<-ctx.Done()
		log.Println("close signal received by walk server.")
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
			if err := recover(); err != nil {
				log.Printf("walking server panic %+v\n", err)
			}
		}()

		s.Shutdown(timeout)
		log.Println("walking server closed.")
	}()

	return s.ListenAndServe()
}

func randomCancel(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		n := rand.Intn(15) + 1
		log.Println("random number is ", n)
		countdown(ctx, n)
		cancel()
	}()
}

func walkHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "walking\n")
	if err != nil {
		log.Fatalf("walk error: %+v", err)
	}
}

func eatHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "eating\n")
	if err != nil {
		log.Fatalf("eat error: %+v", err)
	}
}

func sleepHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "sleeping\n")
	if err != nil {
		log.Fatalf("sleep error: %+v", err)
	}
}

func listenSig(ctx context.Context, cancel context.CancelFunc) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
		log.Printf("cancel sig received.")
		cancel()
		return errors.New("cancel sig received")
	case <-ctx.Done():
		log.Printf("context done.")
		return errors.New("cancel sig received")
	}
}
