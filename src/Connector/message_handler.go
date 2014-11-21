package main

import (
	"github.com/apache/thrift/lib/go/thrift"
	"messageservice"
	"net"
)

var thriftPool = NewClientPool(
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

// func getThriftPool() *ClientPool {

// 	return
// }

func csrAuth(csrInfo *map[interface{}]interface{}) (string, error) {
	// //Auth
	// authClient := authThriftPool.Get()(*authservice.AuthServiceClient)
	// defer authThriftPool.Put(authClient)
	// token := authservice.NewUserToken()
	// userinfo := authClient.getUserinfoByToken()

	// client := thriftPool.Get().(*messageservice.MessageServiceClient)
	// defer thriftPool.Put(client)

	// userId := messageservice.NewUserId()
	// dds := "mobile"
	// userId.Uid = "123"
	// userId.TypeA1 = "csr"
	// userId.ClientType = &dds

	// err := client.Online(userId)
	// return userinfo.uid, err

	return "superfox", nil
}

func customerAuth(customerInfo *map[interface{}]interface{}) (string, error) {
	// authClient := authThriftPool.Get()(*authservice.AuthServiceClient)
	// defer authThriftPool.Put(authClient)
	// token := authservice.NewUserToken()
	// userinfo := authClient.getUserinfoByToken()

	// client := thriftPool.Get().(*messageservice.MessageServiceClient)
	// defer thriftPool.Put(client)

	// userId := messageservice.NewUserId()
	// dds := "mobile"
	// userId.Uid = "123"
	// userId.TypeA1 = "csr"
	// userId.ClientType = &dds

	// err := client.Online(userId)
	// return userinfo.uid, err

	return "superowlf", nil
}

func pipe(message *Message) {
	println(message)
}

func authHandler(message *Message) (string, error) {
	switch message.Cmd {
	case CUSTOMER_AUTH:
		return customerAuth(message.Content)
	case CSR_AUTH:
		return csrAuth(message.Content)
	default:
		panic("传入消息类型不合法!")
	}
}

func messageHandler(message *Message) {
	pipe(message)
}
