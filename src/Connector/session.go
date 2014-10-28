package main

type Session interface {
	subscribe()
	send()
	online()
	offline()
}
