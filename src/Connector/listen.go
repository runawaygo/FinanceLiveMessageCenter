package main

import (
	// "encoding/binary"
	"net"
	// "time"
)

const (
// pingWait = time.Microsecond * 100
// pingWait = time.Second
)

func main() {

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	go h.run()
	go socketIOListen(&h)
	// go statusServicePool.Start()

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		session := NewSession(conn, &h, messageHandler, authHandler, disconnectHandler)
		go session.ReadPump()
		go session.WritePump()
	}
}
