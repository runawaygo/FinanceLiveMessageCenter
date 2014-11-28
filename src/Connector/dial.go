package main

// import (
// 	"fmt"
// 	"net"
// 	"time"
// )

// const (
// // pingWait = time.Microsecond * 100
// // pingWait = time.Second
// )

// func main() {
// 	// conn, err := net.Dial("tcp", "192.168.26.90:8080")
// 	conn, err := net.Dial("tcp", "localhost:8080")

// 	if err != nil {
// 		// handle error
// 		fmt.Println(err)
// 		return
// 	}

// 	session := NewSession(conn, &h, messageHandler, authHandler, disconnectHandler)

// 	go session.ReadPump()
// 	go session.WritePump()

// 	ticker := time.NewTicker(time.Second)
// 	for {
// 		message := Message{Cmd: MESSAGE, Content: &map[interface{}]interface{}{"abc": "superwolf"}}
// 		fmt.Println("superwolf")
// 		session.send <- &message
// 		<-ticker.C

// 		// session.Close()
// 	}
// }
