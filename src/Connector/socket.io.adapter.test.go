package main

import (
	"fmt"
)

import (
	"github.com/garyburd/redigo/redis"
	"github.com/vmihailenco/msgpack"
)

type Message struct {
	Type int64
	Data []string
	Nsp  string
}

type MessageOption struct {
	Except []string
	Rooms  []string
	Flags  map[string]bool
}

func testMsgpack() {
	message := Message{
		Type: 2,
		Data: []string{"superwolf", "fox"},
		Nsp:  "/superowlf",
	}
	message0 := Message{
		Type: 1,
		Nsp:  "/fox",
	}
	msg, err := msgpack.Marshal(message, message0)
	fmt.Println(string(msg))
	var message1 Message
	var message2 Message
	err = msgpack.Unmarshal(msg, &message1, &message2)

	fmt.Println(message1.Nsp)
	fmt.Println(message2.Nsp)
	fmt.Println(message)
	fmt.Println(err)

}

func testMsgpack1() {

}

func main() {
	// testMsgpack1()
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
		return
	}

	// c1, err := redis.Dial("tcp", "127.0.0.1:6379")
	// if err != nil {
	// 	panic(err)
	// 	return
	// }

	psc := redis.PubSubConn{c}
	psc.PSubscribe("socket.io#*")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case redis.PMessage:
			var message = map[string]interface{}{}
			var messageOpt = map[string]interface{}{}
			var messageArray = []interface{}{message, messageOpt}

			msgpack.Unmarshal(v.Data, &messageArray)
			fmt.Println(messageArray)

			// c1.Send("PUBLISH", "socket.io#abc", v.Data)
			// c1.Flush()

		case error:
			return
		}
	}
}

// func BuildMessage(data []byte) []map[string]interface{} {
// 	var messageArray = []map[string]interface{}{map[string]interface{}{}, map[string]interface{}{}}
// 	err := msgpack.Unmarshal(v.Data, &messageArray)
// 	return messageArray, err
// }
