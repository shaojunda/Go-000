package main

import (
	"context"
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
	c := make(chan struct{})
	g, ctx := errgroup.WithContext(context.Background())

	// random cancel
	randomCancel(c)

	// start walking server
	g.Go(func() error {
		return startWalkingServer(c)
	})

	// start eating server
	g.Go(func() error {
		return startEatingServer(c)
	})

	// start sleeping server
	g.Go(func() error {
		return startSleepingServer(c)
	})

	// start bang server
	g.Go(func() error {
		return startBangServer(c)
	})
	listenSig(c)
	go func() {
		for {
			select {
			case <-ctx.Done():
				countdown(10)
				log.Println("All services are closed successfully.")
			}
		}
	}()

	err := g.Wait()

	if err == nil {
		fmt.Println("successfully exit")
	} else {
		fmt.Printf("group error %+v", err)
	}
}

func countdown(seconds int) {
	tick := time.Tick(1 * time.Second)
	for countdown := seconds; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick
	}
}

func startBangServer(c chan struct{}) error {
	m := &http.ServeMux{}
	s := http.Server{Addr: ":8093", Handler: m}
	log.Println("bang server listen on 8093")

	m.HandleFunc("/bang", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "bang!\n")
		if err != nil {
			log.Fatalf("bang error: %+v", err)
		}
		c <- struct{}{}
		timeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		s.Shutdown(timeout)
	})

	return s.ListenAndServe()
}

func startSleepingServer(c chan struct{}) error {
	m := &http.ServeMux{}
	s := http.Server{Addr: ":8092", Handler: m}
	m.HandleFunc("/sleeping", sleepHandler)
	log.Println("sleeping server listen on 8092")

	go func() {
		<-c
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

func startEatingServer(c chan struct{}) error {
	m := &http.ServeMux{}
	s := http.Server{Addr: ":8091", Handler: m}
	m.HandleFunc("/eating", eatHandler)
	log.Println("eat server listen on 8091")

	go func() {
		<-c
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

func startWalkingServer(c chan struct{}) error {
	m := &http.ServeMux{}
	s := http.Server{Addr: ":8090", Handler: m}
	m.HandleFunc("/walking", walkHandler)
	log.Println("walk server listen on 8090")

	go func() {
		<-c
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

func randomCancel(c chan struct{}) {
	go func() {
		n := rand.Intn(15) + 1
		log.Println("random number is ", n)
		countdown(n)
		c <- struct{}{}
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

func listenSig(c chan struct{}) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Print(err)
			}
		}()

		si := <- sig
		log.Println("Got signal :", si)
		if c != nil {
			c <- struct{}{}
		}
	}()
}