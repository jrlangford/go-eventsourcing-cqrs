package secports

import (
	es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/core/domain"
)

// InventoryItemRepository wraps the repository for InventoryItems.
type InventoryItemRepository interface {
	Load(string) (*domain.InventoryItem, error)
	Save(es.AggregateRoot, *int) error
}
