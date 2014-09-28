package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("begin")
	go h.run()
	go loopMessage()

	http.Handle("/", http.FileServer(http.Dir("."))) // <-- note this line
	http.HandleFunc("/ws", serveWs)

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
