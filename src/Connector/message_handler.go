package main

import (
	"github.com/apache/thrift/lib/go/thrift"
	"messageservice"
	"net"
)

var statusServicePool = NewClientPool(
	func() (interface{}, error) {
		trans, err := thrift.NewTSocket(net.JoinHostPort(MESSAGE_SERVICE_CONFIG.Host, MESSAGE_SERVICE_CONFIG.Port))
		if err != nil {
			return nil, err
		}

		err = trans.Open()
		if err != nil {
			return nil, err
		}

		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		client := messageservice.NewMessageServiceClientFactory(trans, protocolFactory)

		return client, nil
	},
	func(resource interface{}) error {
		return resource.(thrift.TTransport).Close()
	},
	MESSAGE_SERVICE_CONFIG.MaxIdle,
	MESSAGE_SERVICE_CONFIG.MaxActive,
)

func pipe(message *Message) error {
	println(message)
	return nil
}

func authHandler(message *Message) (string, string, error) {
	return "superowlf", "mobile", nil
}

func messageHandler(message *Message) error {
	return pipe(message)
}

func disconnectHandler(uid string) {
	// pipe(message)
}
