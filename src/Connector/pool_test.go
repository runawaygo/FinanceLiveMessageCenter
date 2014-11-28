package main

import (

	// "fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ClientPool_Get(t *testing.T) {
	Convey("Normal", t, func() {
		thriftPool := NewClientPool(
			func() (interface{}, error) { return "superwolf", nil },
			func(resource interface{}) error { return nil },
			10,
			10,
		)

		go thriftPool.Start()

		client := thriftPool.Get()
		thriftPool.ClosePool()

		So(client.(string), ShouldEqual, string("superwolf"))
	})
}

func Test_MaxIdle(t *testing.T) {
	Convey("Normal", t, func() {
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
			go func() {
				thriftPool.Put(client)
			}()
		}

		time.Sleep(time.Second * 1)

		So(thriftPool.idleCount, ShouldEqual, 3)
	})
}

func Test_FullLoad(t *testing.T) {
	Convey("Normal", t, func() {
		thriftPool := NewClientPool(
			func() (interface{}, error) { return time.Now(), nil },
			func(resource interface{}) error { return nil },
			1,
			1,
		)
		defer thriftPool.ClosePool()

		go thriftPool.Start()

		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		client := thriftPool.Get()
		go func() {
			<-ticker.C
			thriftPool.Put(client)
		}()

		client2 := thriftPool.Get()

		So(client2.(time.Time).Second(), ShouldEqual, client.(time.Time).Second())
	})

	Convey("Should free resource if idleCount<maxIdleCount", t, func() {
		thriftPool := NewClientPool(
			func() (interface{}, error) { return time.Now(), nil },
			func(resource interface{}) error { return nil },
			0,
			1,
		)
		defer thriftPool.ClosePool()

		go thriftPool.Start()

		client := thriftPool.Get()
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		go func() {
			<-ticker.C
			thriftPool.Put(client)
		}()

		client2 := thriftPool.Get()

		differ := client2.(time.Time).Second() - client.(time.Time).Second()
		So(differ, ShouldEqual, 1)
	})
}
