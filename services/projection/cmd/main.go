package main

import "github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/server"

func main() {
	server := server.New()
	server.Run()
}
