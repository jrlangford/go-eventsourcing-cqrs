package server

import (
	"log"
	"os"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/jrlangford/go-eventsourcing-cqrs/lib/redihash"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/core/usecases"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/primadapters"
	redis "github.com/redis/go-redis/v9"
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
	esdbHost := getEnv("ESDB_HOST", "localhost")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})

	settings, err := esdb.ParseConnectionString("esdb://" + esdbHost + ":2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000")
	if err != nil {
		log.Fatal(err)
	}

	esdbClient, err := esdb.NewClient(settings)
	if err != nil {
		log.Fatal(err)
	}

	projector := usecases.NewProjector(
		redihash.NewHashReadWriter[domain.InventoryItemList](rdb, PUBLIC_INVENTORY_LIST_KEY),
		redihash.NewHashReadWriter[domain.InventoryItemDetails](rdb, PUBLIC_INVENTORY_DETAILS_KEY),
	)

	consumer := primadapters.NewEventConsumer(esdbClient, projector)
	consumer.InitializeSubscription("go-eventsourcing-projector")

	consumer.Run()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
