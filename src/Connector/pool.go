package main

import (
	"container/list"
	"time"
)

type ClientPool struct {
	// Dial is an application supplied function for creating new connections.
	dialFn func() (interface{}, error)

	// Close is an application supplied functoin for closeing connections.
	closeFn func(c interface{}) error

	// TestOnBorrow is an optional application supplied function for checking
	// the health of an idle connection before the connection is used again by
	// the application. Argument t is the time that the connection was returned
	// to the pool. If the function returns an error, then the connection is
	// closed.
	TestOnBorrow func(c interface{}, t time.Time) error

	maxActiveCount int
	maxIdleCount   int

	activeCount int
	idleCount   int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	closed bool

	close   chan bool
	request chan bool
	in      chan interface{}
	out     chan interface{}

	idleList list.List
}

type idleConn struct {
	c interface{}
	t time.Time
}

// New creates a new pool. This function is deprecated. Applications should
// initialize the Pool fields directly as shown in example.
func NewClientPool(dialFn func() (interface{}, error), closeFn func(c interface{}) error, maxIdleCount int, maxActiveCount int) *ClientPool {
	return &ClientPool{
		dialFn:         dialFn,
		closeFn:        closeFn,
		maxIdleCount:   maxIdleCount,
		maxActiveCount: maxActiveCount,
		request:        make(chan bool),
		in:             make(chan interface{}),
		out:            make(chan interface{}),
		close:          make(chan bool),
	}
}

// Get gets a connection. The application must close the returned connection.
// This method always returns a valid connection so that applications can defer
// error handling to the first use of the connection.
func (p *ClientPool) Get() interface{} {
	p.request <- true
	return <-p.out
}

// Put adds conn back to the pool, use forceClose to close the connection forcely
func (p *ClientPool) Put(c interface{}) {
	if p.closed {
		return
	}
	p.in <- c
}

func (p *ClientPool) Start() {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.request:
			if p.activeCount >= p.maxActiveCount {
				p.putToIdle(<-p.in)
				println(p.idleCount)
			}
			p.out <- p.getFromIdle()

		case client := <-p.in:
			p.putToIdle(client)

		case <-p.close:
			break
		}
	}
}

func (p *ClientPool) getFromIdle() interface{} {
	defer func() {
		if x := recover(); x != nil {
			println("run time panic: %v", x)
		}
	}()

	var client interface{}
	var err error

	if p.idleCount == 0 {
		client, err = p.dialFn()
		if err != nil {
			panic(err)
		}
	} else {
		p.idleCount--
		ic := p.idleList.Back()
		p.idleList.Remove(ic)
		client = ic.Value.(idleConn).c
	}

	p.activeCount++
	return client
}

func (p *ClientPool) putToIdle(client interface{}) {
	p.activeCount--
	if p.idleCount >= p.maxIdleCount {
		return
	}
	p.idleList.PushFront(idleConn{t: time.Now(), c: client})
	p.idleCount++
}

// ClosePool releases the resources used by the pool.
func (p *ClientPool) ClosePool() {
	if p.closed {
		return
	}

	p.closed = true
	p.close <- true

	for e := p.idleList.Front(); e != nil; e = e.Next() {
		err := p.closeFn(e.Value.(idleConn).c)
		println(err)
	}

	p.idleCount = 0
	p.activeCount = 0
	p.idleList.Init()

}
