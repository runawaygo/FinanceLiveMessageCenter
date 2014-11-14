package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// import "runtime/debug"
import (
	"github.com/vmihailenco/msgpack"
)

const (
	// Time allowed to write a message to the peer.
	// pingWait   = 10 * time.Second
	// pingWait   = time.Second * 10
	pingPeriod = (pingWait * 9) / 10
)

var authMap = map[uint16]string{CSR_AUTH: "csr", CUSTOMER_AUTH: "customer"}

type Session struct {
	Reader         *bufio.Reader
	Conn           net.Conn
	send           chan Message
	hub            IHub
	uid            string
	messageHandler func(message *Message)
	authHandler    func(message *Message) (string, error)
}

func NewSession(
	conn net.Conn,
	h IHub,
	messageHandler func(message *Message),
	authHandler func(message *Message) (string, error),
) *Session {

	s := new(Session)
	s.Conn = conn
	s.Reader = bufio.NewReader(conn)
	s.send = make(chan Message)
	s.hub = h
	s.uid = ""
	s.messageHandler = messageHandler
	s.authHandler = authHandler
	return s
}

func (session *Session) deferFunc() {
	session.Conn.Close()
	session.hub.Unregister(session)

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
	if _, err := session.Reader.ReadBytes('$'); err != nil {
		panic(err)
	}

	char, err1 := session.Reader.ReadByte()
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
	session.sendCmd(message.Cmd)

	objBytes, err := msgpack.Marshal(*message.Content)
	if err != nil {
		panic(err)
	}

	data := []byte{}
	data = append(data, convertToByte(uint16(len(objBytes)))...)
	data = append(data, objBytes...)

	if _, err1 := session.Conn.Write(data); err1 != nil {
		panic(err1)
	}
}

func (session *Session) sendCmd(cmd uint16) {
	data := []byte{'$', '$'}
	data = append(data, convertToByte(cmd)...)

	if _, err1 := session.Conn.Write(data); err1 != nil {
		panic(err1)
	}
}

func (session *Session) auth(nsp string, message *Message) {
	uid, err := session.authHandler(message)
	if err != nil {
		session.sendCmd(AUTH_FAILED)
		return
	}
	session.uid = uid
	h.register <- session
	h.Join(nsp, session.uid, uid)
}

func (session *Session) WritePump() {
	defer session.deferFunc()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	session.Conn.SetWriteDeadline(time.Now().Add(pingWait))

	for {
		fmt.Println("write pump")
		select {
		case message, ok := <-session.send:
			if !ok {
				session.sendCmd(INTERNAL_ERROR)
				return
			}
			session.sendMessage(&message)

		case <-ticker.C:
			session.sendCmd(PUMP)
			session.Conn.SetWriteDeadline(time.Now().Add(pingWait))
		}
	}
}

func (session *Session) Close() {
	defer session.deferFunc()
}

func (session *Session) ReadPump() {
	defer session.deferFunc()

	session.Conn.SetReadDeadline(time.Now().Add(pingWait))

	for {
		fmt.Println("read pump")
		cmd := session.readCmd()
		session.Conn.SetReadDeadline(time.Now().Add(pingWait))

		if cmd >= 10000 {
			dict := map[interface{}]interface{}{}
			session.readMessage(&dict)
			message := &Message{Cmd: cmd, Content: &dict}

			if nsp, ok := authMap[cmd]; ok {
				session.auth(nsp, message)
				continue
			}

			session.messageHandler(message)
		}
	}
}
