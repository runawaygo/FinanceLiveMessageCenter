package main

import (
	"bufio"
	// "encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	reader := bufio.NewReader(conn)
	for {

		data := readMessage(reader)
		fmt.Println(data)

		<-ticker.C
		fmt.Println(time.Now())
	}
}
