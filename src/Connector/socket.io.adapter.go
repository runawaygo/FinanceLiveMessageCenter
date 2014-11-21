package main

import (
	"fmt"
)

import (
	"github.com/garyburd/redigo/redis"
	"github.com/vmihailenco/msgpack"
)

type SocketIOMessage struct {
	Nsp  string
	Type string
	Data map[interface{}]interface{}
	Room string
}

var cmdMap = map[string]uint16{
	"message": MESSAGE,
}

func socketIOListen(h IHub) {
	c, err := redis.Dial("tcp", REDIS_CONFIG.Host)
	if err != nil {
		panic(err)
		return
	}

	psc := redis.PubSubConn{c}
	psc.PSubscribe("socket.io#*")
	for {
		switch v := psc.Receive().(type) {
		case redis.PMessage:
			var messageData = map[string]interface{}{}
			var messageOpt = map[string]interface{}{}
			var messageArray = []interface{}{messageData, messageOpt}

			msgpack.Unmarshal(v.Data, &messageArray)
			message := &SocketIOMessage{
				Nsp:  messageData["nsp"].(string),
				Room: messageOpt["rooms"].([]interface{})[0].(string),
				Type: messageData["data"].([]interface{})[0].(string),
				Data: messageData["data"].([]interface{})[1].(map[interface{}]interface{}),
			}

			fmt.Println(message)

			h.Broadcast(&Broadcast{
				Nsp:  message.Nsp,
				Room: message.Room,
				Message: Message{
					Cmd:     cmdMap[message.Type],
					Content: &message.Data,
				},
			})

		case error:
			return
		}
	}
}
