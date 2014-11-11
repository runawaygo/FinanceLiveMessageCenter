package main

type Message struct {
	Cmd     uint16
	Content *map[interface{}]interface{}
}

type Broadcast struct {
	Nsp     string
	Room    string
	Message Message
}
