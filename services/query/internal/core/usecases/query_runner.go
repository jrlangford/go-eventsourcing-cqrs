package usecases

import (
	"log"

	"github.com/jrlangford/go-eventsourcing-cqrs/lib/query"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/secports"
	redis "github.com/redis/go-redis/v9"
)

// A QueryRunner is a façade to the query handlers.
type QueryRunner struct {
	rdb        *redis.Client
	dispatcher *query.InMemoryQueryDispatcher
}

// NewQueryRunner is self-describing.
func NewQueryRunner(
	rdb *redis.Client,
	itemListRepo secports.InventoryItemListRepo,
	itemDetailsRepo secports.InventoryItemDetailsRepo,
) *QueryRunner {

	invQueryHandler := newInventoryQueryHandler(rdb, itemListRepo, itemDetailsRepo)

	dispatcher := query.NewInMemoryQueryDispatcher()

	err := dispatcher.RegisterHandler(invQueryHandler,
		&domain.GetInventoryItems{},
		&domain.GetInventoryItemDetails{},
	)
	if err != nil {
		log.Fatal(err)
	}
	return &QueryRunner{
		rdb:        rdb,
		dispatcher: dispatcher,
	}
}

// Run provides a query response corresponding to the provided query message.
func (qr *QueryRunner) Run(message query.QueryMessage) (interface{}, error) {
	log.Printf("%+v", message)
	return qr.dispatcher.Dispatch(message)
}
