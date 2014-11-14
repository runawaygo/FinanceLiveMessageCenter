package main

import (
	"sync"
	"time"
)

type ClientPool struct {
	sync.Pool
	// Dial is an application supplied function for creating new connections.
	Dial func() (interface{}, error)

	// Close is an application supplied functoin for closeing connections.
	Close func(c interface{}) error

	// TestOnBorrow is an optional application supplied function for checking
	// the health of an idle connection before the connection is used again by
	// the application. Argument t is the time that the connection was returned
	// to the pool. If the function returns an error, then the connection is
	// closed.
	TestOnBorrow func(c interface{}, t time.Time) error

	maxCount     int
	maxIdleCount int

	activeCount int
	idleCount   int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	closed  bool
	request chan bool
	push    chan idleConn
	pop     chan idleConn
	isBack  chan int
}

type idleConn struct {
	c interface{}
	t time.Time
}

// New creates a new pool. This function is deprecated. Applications should
// initialize the Pool fields directly as shown in example.
func New(dialFn func() (interface{}, error), closeFn func(c interface{}) error, maxIdle int) *Pool {
	return &Pool{Dial: dialFn, Close: closeFn, MaxIdle: maxIdle}
}

// Get gets a connection. The application must close the returned connection.
// This method always returns a valid connection so that applications can defer
// error handling to the first use of the connection.
func (p *Pool) GetClient() interface{} {
	p.request <- true
	return <-p.pop
}

// Put adds conn back to the pool, use forceClose to close the connection forcely
func (p *Pool) PutClient(c interface{}) error {
	p.push <- c
}

func (p *Pool) Start() {
	var client idleConn
	for {
		fmt.Println("write pump")
		select {
		case <-p.request:
			var client interface{}
			if p.activeCount >= p.maxCount {
				client = <-p.push
			} else {
				if p.idleCount > 0 {
					p.idleCount--
				}
				p.activeCount++
				client = p.Get().(idleConn).c
			}
			p.pop <- client

		case client := <-p.push:
			p.activeCount--
			if p.idleCount >= p.maxIdleCount {
				continue
			}
			p.Put(idleConn{c: client, t: time.Now()})
			p.idleCount++
		}
	}
}

// Relaase releases the resources used by the pool.
func (p *Pool) Release() error {
	p.mu.Lock()
	idle := p.idle
	p.idle.Init()
	p.closed = true
	p.active -= idle.Len()
	p.mu.Unlock()
	for e := idle.Front(); e != nil; e = e.Next() {
		p.Close(e.Value.(idleConn).c)
	}
	return nil
}
