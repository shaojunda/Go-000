package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	c := make(chan string)
	flag := rand.Intn(100)
	rd := bufio.NewReader(conn)
	go echo(conn, c)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			log.Printf("read error: %v\n", err)
			return
		}
		msg := string(line)
		fmt.Printf("%d say: %v\n", flag, msg)
		c <- msg
	}
}

func echo(conn net.Conn, c <-chan string) {
	wr := bufio.NewWriter(conn)
	for msg := range c {
		if msg == "q!" {
			wr.WriteString(strings.Join([]string{"echo~: ", "bye~", "\n"}, ""))
			wr.Flush()
			conn.Close()
		}
		wr.WriteString(strings.Join([]string{"echo~: ", msg, "\n"}, ""))
		wr.Flush()
	}
}
