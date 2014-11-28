package main

// import (
// 	// "fmt"
// 	// . "github.com/gbbr/mocks"
// 	"github.com/maraino/go-mock"
// 	// . "github.com/smartystreets/goconvey/convey"
// 	// "net"
// 	// "time"
// )

// type funcStub struct {
// 	mock.Mock
// }

// func (s *funcStub) abc() {
// 	s.Called()

// }

// func (s *funcStub) authHandlerMock(message string) (string, error) {
// 	ret := s.Called(message)
// 	return ret.String(0), ret.Error(1)
// }

// func main() {
// 	fStub := &funcStub{}
// 	// fStub.authHandlerMock("superwolf")
// 	fStub.When("abc").Return()
// 	fStub.abc()
// 	// c1, c2 := Pipe(
// 	// 	&Conn{RAddr: "1.1.1.1:123"},
// 	// 	&Conn{LAddr: "127.0.0.1:12", RAddr: "2.2.2.2:456"},
// 	// )

// 	// session := NewSession(c1, &h, messageHandler, authHandlerMock, disconnectHandler)
// 	// go session.ReadPump()
// 	// go session.WritePump()

// 	// content := map[interface{}]interface{}{}

// 	// message := &Message{
// 	// 	Cmd:     1100,
// 	// 	Content: &content,
// 	// }

// 	// c2.Write(message.toBytes())

// 	// time.Sleep(time.Microsecond * 20)

// }
