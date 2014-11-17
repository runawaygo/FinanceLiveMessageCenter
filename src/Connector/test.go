package main

import (
	// "fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"messageservice"
	"net"
	"sync"
)

type abc struct {
	sync.Pool
}

func main() {
	thriftPool := NewClientPool(
		func() (interface{}, error) { return "superwolf", nil },
		func(resource interface{}) error { return nil },
		3,
		30,
	)
	defer thriftPool.ClosePool()

	go thriftPool.Start()

	for i := 0; i < 30; i++ {
		client := thriftPool.Get()
		// go func() {
		// 	thriftPool.Put(client)
		// }()
		println(thriftPool.activeCount)
		println(thriftPool.idleCount)
	}
}

func main1() {
	dd := &abc{}
	dd.Pool.Put("hehe")
	aaa := dd.Pool.Get()
	println("hehe")
	println(dd)
	println(aaa.(string))

	return
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
