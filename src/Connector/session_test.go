package main

import (
	"errors"
	. "github.com/gbbr/mocks"
	"github.com/maraino/go-mock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type hubStub struct {
	mock.Mock
}

func (s *hubStub) Broadcast(broadcast *Broadcast) {
	s.Called(broadcast)
}

func (s *hubStub) Register(session ISession) {
	s.Called(session)
}

func (s *hubStub) Unregister(session ISession) {
	s.Called(session)
}

type funcStub struct {
	mock.Mock
}

func (s *funcStub) authHandlerMock(message *Message) (string, string, error) {
	ret := s.Called(message)
	return ret.String(0), ret.String(1), ret.Error(2)
}

func (s *funcStub) messageHandlerMock(message *Message) error {
	ret := s.Called(message)
	return ret.Error(0)
}

func Test_Session(t *testing.T) {

	Convey("Session", t, func() {
		c1, c2 := Pipe(&Conn{}, &Conn{})

		Convey("Pump", func() {
			fStub := &funcStub{}

			session := NewSession(c1, &h, fStub.messageHandlerMock, fStub.authHandlerMock, disconnectHandler)

			go session.WritePump()
			time.Sleep(time.Microsecond * 1200)

			resultBytes := make([]byte, 1000)
			c2.Read(resultBytes)
			So(resultBytes[3], ShouldEqual, 3)

		})

		Convey("In", func() {
			message := &Message{
				Cmd: 1100,
				Content: &map[interface{}]interface{}{
					"username": "superwolf",
					"password": "holyshit",
				},
			}

			Convey("AuthSuccess", func() {
				fStub := &funcStub{}
				session := NewSession(c1, &h, fStub.messageHandlerMock, fStub.authHandlerMock, disconnectHandler)

				go session.ReadPump()
				go session.WritePump()

				fStub.When("authHandlerMock", message).Return("superwolf", "mobile", nil).Times(1)

				c2.Write(message.toBytes())

				time.Sleep(time.Microsecond * 20)

				if ok, err := fStub.Verify(); !ok {
					So(err, ShouldBeNil)
				}
			})

			Convey("AuthFailure", func() {
				fStub := &funcStub{}
				session := NewSession(c1, &h, fStub.messageHandlerMock, fStub.authHandlerMock, disconnectHandler)

				go session.ReadPump()
				go session.WritePump()

				fStub.When("authHandlerMock", message).Return("", "", errors.New("superwolf")).Times(1)

				c2.Write(message.toBytes())

				time.Sleep(time.Microsecond * 20)

				resultBytes := make([]byte, 4)

				c2.Read(resultBytes)

				So(resultBytes[2], ShouldEqual, 3)
				So(resultBytes[3], ShouldEqual, 232)

				if ok, err := fStub.Verify(); !ok {
					So(err, ShouldBeNil)
				}
			})

			Convey("ReadMessage", func() {
				message := &Message{
					Cmd: MESSAGE,
					Nsp: "",
					Uid: "",
					Content: &map[interface{}]interface{}{
						"username": "superwolf",
						"password": "holyshit",
					},
				}
				fStub := &funcStub{}
				session := NewSession(c1, &h, fStub.messageHandlerMock, fStub.authHandlerMock, disconnectHandler)

				go session.ReadPump()
				go session.WritePump()

				fStub.When("messageHandlerMock", message).Return(nil).Times(1)

				c2.Write(message.toBytes())

				time.Sleep(time.Microsecond * 20)

				if ok, err := fStub.Verify(); !ok {
					So(err, ShouldBeNil)
				}
			})

		})

		Convey("Out", func() {
			message := &Message{
				Cmd: MESSAGE,
				Content: &map[interface{}]interface{}{
					"type":    33,
					"name":    "superwolf",
					"content": "hehe",
				},
			}

			Convey("Hub Broadcast", func() {
				fStub := &funcStub{}

				session := NewSession(c1, &h, fStub.messageHandlerMock, fStub.authHandlerMock, disconnectHandler)

				go session.ReadPump()
				go session.WritePump()

				session.Send(message)

				resultBytes := make([]byte, 1000)
				c2.Read(resultBytes)

				time.Sleep(time.Microsecond * 20)

				cmd := convertToUint16([]byte{resultBytes[2], resultBytes[3]})
				So(cmd, ShouldEqual, MESSAGE)

				if ok, err := fStub.Verify(); !ok {
					So(err, ShouldBeNil)
				}
			})
		})
	})
}
