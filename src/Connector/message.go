package main

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
)

type Message struct {
	Cmd     uint16
	Nsp     string
	Uid     string
	Content *map[interface{}]interface{}
}

type Broadcast struct {
	Nsp     string
	Room    string
	Message *Message
}

func (message *Message) toBytes() []byte {
	data := []byte{'$', '$'}
	data = append(data, convertToByte(message.Cmd)...)

	if message.Content == nil || len(*message.Content) == 0 {
		return data
	}

	objBytes, err := msgpack.Marshal(*message.Content)
	if err != nil {
		fmt.Printf("run time panic: %v", err)
		return nil
	}

	data = append(data, convertToByte(uint16(len(objBytes)))...)
	data = append(data, objBytes...)
	return data
}
