// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type IHub interface {
	Broadcast(broadcast *Broadcast)
	Register(session ISession)
	Unregister(session ISession)
}

type nsp struct {
	sessions map[string]ISession
	rooms    map[string]map[string]bool
}

type hub struct {
	// Registered connections.
	nsps map[string]*nsp

	broadcast chan *Broadcast

	register chan ISession

	unregister chan ISession
}

var h = hub{
	broadcast:  make(chan *Broadcast),
	register:   make(chan ISession),
	unregister: make(chan ISession),
	nsps:       make(map[string]*nsp),
}

func (h *hub) Join(nspName string, room string, uid string) {
	if _, ok := h.nsps[nspName].sessions[uid]; !ok {
		return
	}

	if _, ok := h.nsps[nspName].rooms[room]; !ok {
		h.nsps[nspName].rooms[room] = make(map[string]bool)
	}

	h.nsps[nspName].rooms[room][uid] = true
}

func (h *hub) Left(nspName string, room string, uid string) {
	if _, ok := h.nsps[nspName].sessions[uid]; !ok {
		return
	}

	if _, ok := h.nsps[nspName].rooms[room]; !ok {
		return
	}

	delete(h.nsps[nspName].rooms[room], uid)
}

func (h *hub) Register(session ISession) {
	nspName := session.GetNsp()
	if _, ok := h.nsps[nspName]; !ok {
		h.nsps[nspName] = &nsp{
			sessions: make(map[string]ISession),
			rooms:    make(map[string]map[string]bool),
		}
	}
	h.register <- session

	uid := session.GetUid()
	h.Join(nspName, uid, uid)
}

func (h *hub) Unregister(session ISession) {
	nspName := session.GetNsp()
	uid := session.GetUid()

	h.Left(nspName, uid, uid)
	h.unregister <- session
}

func (h *hub) Broadcast(broadcast *Broadcast) {
	h.broadcast <- broadcast
}

func (h *hub) BroadcastWithArgs(nsp string, room string, message *Message) {
	h.Broadcast(&Broadcast{
		Nsp:     nsp,
		Room:    room,
		Message: message,
	})
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			sessionId := c.GetUid()
			nspName := c.GetNsp()

			h.nsps[nspName].sessions[sessionId] = c

		case c := <-h.unregister:
			sessionId := c.GetUid()
			nspName := c.GetNsp()

			if _, ok := h.nsps[nspName].sessions[sessionId]; ok {
				delete(h.nsps[nspName].sessions, sessionId)
			}

		case b := <-h.broadcast:
			if _, ok := h.nsps[b.Nsp].rooms[b.Room]; !ok {
				continue
			}
			sessions := h.nsps[b.Nsp].sessions
			for sessionId := range h.nsps[b.Nsp].rooms[b.Room] {
				if session, ok := sessions[sessionId]; ok {
					session.Send(b.Message)
					continue
				}

				delete(h.nsps[b.Nsp].rooms[b.Room], sessionId)
			}
		}
	}
}
