package server

import (
	"log"
	"net"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/core/usecases"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/primadapters"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/primadapters/pb"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/secadapters"
	"google.golang.org/grpc"
)

// A server manages the server's resources.
type server struct {
}

// New constructs a new server.
func New() *server {
	return &server{}
}

// Run initilizes and runs the server.
func (srv *server) Run() {

	settings, err := esdb.ParseConnectionString("esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000")

	if err != nil {
		panic(err)
	}

	esClient, err := esdb.NewClient(settings)
	if err != nil {
		log.Fatal(err)
	}

	invItemsRepo, err := secadapters.NewInventoryItemRepo(esClient)
	if err != nil {
		log.Fatal(err)
	}

	executor := usecases.NewExecutor(invItemsRepo)

	commandServer := primadapters.NewCommandServer(executor)

	server := grpc.NewServer()
	pb.RegisterCommandServer(server, commandServer)

	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server is running on port 50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
