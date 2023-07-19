package usecases

import (
	"log"

	es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/core/secports"
)

// An Executor is a façade for the command handlers.
type Executor struct {
	dispatcher *es.InMemoryDispatcher
}

// NewExecutor is self-describing.
func NewExecutor(invItemsRepo secports.InventoryItemRepository) *Executor {

	inventoryCommandHandler := newInventoryCommandHandler(invItemsRepo)

	dispatcher := es.NewInMemoryDispatcher()

	err := dispatcher.RegisterHandler(inventoryCommandHandler,
		&domain.CreateInventoryItem{},
		&domain.DeactivateInventoryItem{},
		&domain.RenameInventoryItem{},
		&domain.CheckInItemsToInventory{},
		&domain.RemoveItemsFromInventory{},
	)
	if err != nil {
		log.Fatal(err)
	}
	return &Executor{
		dispatcher: dispatcher,
	}
}

// Execute executes a command message using its corresponding handler, if
// present.
func (ex *Executor) Execute(command es.CommandMessage) error {
	return ex.dispatcher.Dispatch(command)
}
