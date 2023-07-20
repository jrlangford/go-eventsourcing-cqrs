package server

import (
	"log"
	"net"
	"os"

	"github.com/jrlangford/go-eventsourcing-cqrs/lib/redihash"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/secports"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/usecases"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/primadapters"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/primadapters/pb"
	redis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

const PUBLIC_INVENTORY_NAMESPACE = "public:inventory"
const PUBLIC_INVENTORY_LIST_KEY = PUBLIC_INVENTORY_NAMESPACE + ":list"
const PUBLIC_INVENTORY_DETAILS_KEY = PUBLIC_INVENTORY_NAMESPACE + ":details"

// A server manages the server's resources.
type server struct {
}

// New constructs a new server.
func New() *server {
	return &server{}
}

// Run initilizes and runs the server.
func (srv *server) Run() {

	redisHost := getEnv("REDIS_HOST", "localhost")
	queryListenAddress := getEnv("QUERY_LISTEN_ADDRESS", "localhost")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})

	qr := usecases.NewQueryRunner(
		rdb,
		redihash.NewHashReader[secports.InventoryItemListDto](rdb, PUBLIC_INVENTORY_LIST_KEY),
		redihash.NewHashReader[secports.InventoryItemDetailsDto](rdb, PUBLIC_INVENTORY_DETAILS_KEY),
	)

	queryServer := primadapters.NewQueryServer(qr)

	server := grpc.NewServer()
	pb.RegisterQueryServer(server, queryServer)

	listener, err := net.Listen("tcp", queryListenAddress + ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server is running on port 50052")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
