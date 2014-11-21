package main

import (
	// "fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"messageservice"
	"net"
	"sync"
	"time"
)

type abc struct {
	sync.Pool
}

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

func main() {
	go thriftPool.Start()
	client := thriftPool.Get().(*messageservice.MessageServiceClient)
	defer func() {
		thriftPool.Put(client)

		time.Sleep(time.Second * 2)

		println(thriftPool.idleCount)
		println(thriftPool.idleList.Len())
	}()

	userId := messageservice.NewUserId()
	dds := "mobile"
	userId.Uid = "123"
	userId.TypeA1 = "csr"
	userId.ClientType = &dds

	err := client.Online(userId)
	if err != nil {
		panic(err)
	}

	userIdCollection := messageservice.UserIdCollection{userId}
	results, err := client.GetUserOnlineStatus(userIdCollection)
	if err != nil {
		panic(err)
	}
	println(len(results))
	println(results[0])

	err = client.Offline(userId)
	if err != nil {
		panic(err)
	}

	results, err = client.GetUserOnlineStatus(userIdCollection)
	if err != nil {
		panic(err)
	}
	println(len(results))
	println(results[0])

	println(thriftPool.idleCount)
	println(thriftPool.idleList.Len())
	println(thriftPool.activeCount)
}

func main1() {
	// dd := &abc{}
	// dd.Pool.Put("hehe")
	// aaa := dd.Pool.Get()
	// println("hehe")
	// println(dd)
	// println(aaa.(string))

	// return
	trans, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "10001"))
	if err != nil {
		panic(err)
	}

	err = trans.Open()
	if err != nil {
		panic(err)
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := messageservice.NewMessageServiceClientFactory(trans, protocolFactory)

	userId := messageservice.NewUserId()
	dds := "mobile"
	userId.Uid = "123"
	userId.TypeA1 = "csr"
	userId.ClientType = &dds

	err = client.Online(userId)
	if err != nil {
		panic(err)
	}

	userIdCollection := messageservice.UserIdCollection{userId}
	results, err := client.GetUserOnlineStatus(userIdCollection)
	if err != nil {
		panic(err)
	}
	println(len(results))
	println(results[0])

	err = client.Offline(userId)
	if err != nil {
		panic(err)
	}

	results, err = client.GetUserOnlineStatus(userIdCollection)
	if err != nil {
		panic(err)
	}
	println(len(results))
	println(results[0])

}
