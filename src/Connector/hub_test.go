package main

import (

	// "fmt"
	"testing"
	"time"

	"github.com/maraino/go-mock"
	. "github.com/smartystreets/goconvey/convey"
)

type sessionStub struct {
	mock.Mock
}

func (s *sessionStub) GetUid() string {
	ret := s.Called()
	return ret.String(0)
}

func (s *sessionStub) GetNsp() string {
	ret := s.Called()
	return ret.String(0)
}

func (s *sessionStub) Send(m *Message) {
	s.Called(m)
}

func genHub() *hub {
	return &hub{
		broadcast:  make(chan *Broadcast),
		register:   make(chan ISession),
		unregister: make(chan ISession),
		nsps:       make(map[string]*nsp),
	}
}

func Test_Hub(t *testing.T) {
	message := &Message{
		Cmd: MESSAGE,
		Content: &map[interface{}]interface{}{
			"type":    33,
			"name":    "superwolf",
			"content": "hehe",
		},
	}

	Convey("Normal", t, func() {
		s := &sessionStub{}
		s.When("GetUid").Return("superwolf")
		s.When("GetNsp").Return("csr")

		Convey("Hub status", func() {
			h := genHub()
			So(h.nsps["csr"], ShouldBeNil)
		})

		Convey("Broadcast after join", func() {
			h := genHub()
			go h.run()
			s.When("Send", message).Times(1)

			//注册Session
			h.Register(s)
			So(h.nsps["csr"].sessions["superwolf"], ShouldEqual, s)

			//Join room
			h.Join("csr", "helo", "superwolf")
			So(len(h.nsps["csr"].rooms["helo"]), ShouldEqual, 1)

			//Send Broadcast
			h.BroadcastWithArgs("csr", "helo", message)
			time.Sleep(time.Second)

			if ok, err := s.Verify(); !ok {
				So(err, ShouldBeNil)
			}
		})

		Convey("Broadcast after left", func() {
			h := genHub()
			go h.run()
			s.When("Send", message).Times(0)

			//注册Session
			h.Register(s)
			time.Sleep(time.Microsecond * 20)
			h.Join("csr", "helo", "superwolf")
			h.Left("csr", "helo", "superwolf")

			h.BroadcastWithArgs("csr", "helo", message)
			time.Sleep(time.Microsecond * 20)
			if ok, err := s.Verify(); !ok {
				So(err, ShouldBeNil)
			}
		})

		Convey("Broadcast after unregister", func() {
			h := genHub()
			go h.run()
			s.When("Send", message).Times(0)

			//注册Session
			h.Register(s)
			time.Sleep(time.Microsecond * 20)
			h.Join("csr", "helo", "superwolf")

			h.Unregister(s)
			time.Sleep(time.Microsecond * 20)

			h.BroadcastWithArgs("csr", "helo", message)
			time.Sleep(time.Microsecond * 20)

			_, err := s.Verify()
			So(err, ShouldBeNil)
			So(h.nsps["csr"].rooms["helo"]["superwolf"], ShouldBeFalse)
		})
	})
}
