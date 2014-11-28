package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

// import "runtime/debug"
import (
	"github.com/vmihailenco/msgpack"
)

const (
	// Time allowed to write a message to the peer.
	pingWait   = time.Second
	pingPeriod = (pingWait * 9) / 10
)

type ISession interface {
	Send(message *Message)
	GetUid() string
	GetNsp() string
}

type Session struct {
	*bufio.Reader
	net.Conn

	send       chan *Message
	hub        IHub
	isAuthed   bool
	uid        string
	clientType string
	nsp        string

	once sync.Once

	messageHandler    func(message *Message) error
	authHandler       func(message *Message) (string, string, error)
	disconnectHandler func(uid string)
}

func NewSession(
	conn net.Conn,
	h IHub,
	messageHandler func(message *Message) error,
	authHandler func(message *Message) (string, string, error),
	disconnectHandler func(uid string),
) *Session {

	s := new(Session)
	s.Conn = conn
	s.Reader = bufio.NewReader(conn)
	s.send = make(chan *Message)
	s.hub = h

	s.isAuthed = false
	s.uid = ""
	s.clientType = ""
	s.nsp = ""

	s.messageHandler = messageHandler
	s.authHandler = authHandler
	s.disconnectHandler = disconnectHandler
	return s
}

func (session *Session) GetUid() string {
	return session.uid
}

func (session *Session) GetNsp() string {
	return session.nsp
}

func (session *Session) deferFunc() {
	session.once.Do(func() {
		if session.isAuthed {
			session.hub.Unregister(session)
		}

		close(session.send)

		session.Conn.Close()
	})

	if x := recover(); x != nil {
		fmt.Printf("run time panic: %v", x)
	}
}

func (session *Session) readUint16() uint16 {
	cmdBytes := make([]byte, 2)
	if _, err := session.Reader.Read(cmdBytes); err != nil {
		panic(err)
	}

	return convertToUint16(cmdBytes)
}

func (session *Session) skipBeginFlag() {
	if _, err := session.ReadBytes('$'); err != nil {
		panic(err)
	}

	char, err1 := session.ReadByte()
	if err1 != nil {
		panic(err1)
	}

	if char != '$' {
		session.skipBeginFlag()
	}
	return
}

func (session *Session) readCmd() uint16 {
	session.skipBeginFlag()
	return session.readUint16()
}

func (session *Session) readMessage(dict *map[interface{}]interface{}) {
	length := session.readUint16()

	messageBytes := make([]byte, length)
	if _, err := session.Reader.Read(messageBytes); err != nil {
		panic(err)
	}

	msgpack.Unmarshal(messageBytes, dict)
}

func (session *Session) sendMessage(message *Message) {
	_, err1 := session.Write(message.toBytes())

	if err1 != nil {
		panic(err1)
	}
}

func (session *Session) sendCmd(cmd uint16) {
	data := []byte{'$', '$'}
	data = append(data, convertToByte(cmd)...)

	if _, err1 := session.Write(data); err1 != nil {
		panic(err1)
	}
}

func (session *Session) auth(message *Message) {
	uid, clientType, err := session.authHandler(message)
	if err != nil {
		session.sendCmd(AUTH_FAILED)
		return
	}

	session.isAuthed = true

	session.nsp = string(message.Cmd)
	session.uid = uid
	session.clientType = clientType

	session.hub.Register(session)
}

func (session *Session) WritePump() {
	defer session.deferFunc()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	session.SetWriteDeadline(time.Now().Add(pingWait))

	for {
		fmt.Println("write pump")
		select {
		case message, ok := <-session.send:
			if !ok {
				session.sendCmd(INTERNAL_ERROR)
				return
			}
			session.sendMessage(message)

		case <-ticker.C:
			session.sendCmd(PUMP)
			session.SetWriteDeadline(time.Now().Add(pingWait))
		}
	}
}

func (session *Session) Close() {
	defer session.deferFunc()
}

func (session *Session) ReadPump() {
	defer session.deferFunc()

	session.SetReadDeadline(time.Now().Add(pingWait))

	for {
		fmt.Println("read pump")
		cmd := session.readCmd()

		println(cmd)
		session.SetReadDeadline(time.Now().Add(pingWait))

		switch {
		case cmd == CLOSE:
			return

		case cmd < 1000:
			continue

		case 1100 <= cmd && cmd < 1300:
			dict := map[interface{}]interface{}{}
			session.readMessage(&dict)
			message := &Message{Cmd: cmd, Content: &dict}

			session.auth(message)
		default:
			dict := map[interface{}]interface{}{}
			session.readMessage(&dict)
			message := &Message{Cmd: cmd, Nsp: session.nsp, Uid: session.uid, Content: &dict}

			session.messageHandler(message)
		}
		println("superowlf")
	}
}

func (session *Session) Send(message *Message) {
	session.send <- message
}
