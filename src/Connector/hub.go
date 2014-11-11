// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type IHub interface {
	Broadcast(broadcast *Broadcast)
	Register(session *Session)
	Unregister(session *Session)
}

type hub struct {
	// Registered connections.
	connections map[*Session]string

	sessions map[string]*Session

	nsp map[string]map[string]map[*Session]bool

	broadcast chan *Broadcast

	register chan *Session

	unregister chan *Session
}

var h = hub{
	broadcast:   make(chan *Broadcast),
	register:    make(chan *Session),
	unregister:  make(chan *Session),
	connections: make(map[*Session]string),
	sessions:    make(map[string]*Session),
	nsp:         make(map[string]map[string]map[*Session]bool),
}

func (h *hub) Join(nsp string, room string, uid string) {
	if _, ok := h.sessions[uid]; !ok {
		return
	}

	if _, ok := h.nsp[nsp]; !ok {
		h.nsp[nsp] = make(map[string]map[*Session]bool)
	}

	if _, ok := h.nsp[nsp][room]; !ok {
		h.nsp[nsp][room] = make(map[*Session]bool)
	}

	h.nsp[nsp][room][h.sessions[uid]] = true
}

func (h *hub) Left(nsp string, room string, uid string) {
	if _, ok := h.sessions[uid]; !ok {
		return
	}

	if _, ok := h.nsp[nsp][room]; !ok {
		return
	}

	delete(h.nsp[nsp][room], h.sessions[uid])
}

func (h *hub) Register(session *Session) {
	h.register <- session
}

func (h *hub) Unregister(session *Session) {
	h.unregister <- session
}

func (h *hub) Broadcast(broadcast *Broadcast) {
	h.broadcast <- broadcast
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = c.uid
			h.sessions[c.uid] = c

		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
			if _, ok := h.sessions[c.uid]; ok {
				delete(h.sessions, c.uid)
			}

		case b := <-h.broadcast:
			if _, ok := h.nsp[b.Nsp][b.Room]; !ok {
				continue
			}

			for c := range h.nsp[b.Nsp][b.Room] {
				select {
				case c.send <- b.Message:
				default:
					close(c.send)
					delete(h.nsp[b.Nsp][b.Room], c)
				}
			}
		}
	}
}
