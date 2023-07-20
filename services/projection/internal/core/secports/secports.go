package secports

import (
	"context"

	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/core/domain"
)

// ProjectionRepo wraps the methods for the repository of the provided type.
type ProjectionRepo[T any] interface {
	Write(ctx context.Context, dataKey string, val *T) error
	Update(ctx context.Context, dataKey string, updater func(*T)) error
	Del(ctx context.Context, dataKey string) error
}

// InventoryItemListRepo is an alias for an InventoryItemList projection repo.
type InventoryItemListRepo = ProjectionRepo[domain.InventoryItemList]

// InventoryItemDetailsRepo is an alias for an InventoryItemDetails projection repo.
type InventoryItemDetailsRepo = ProjectionRepo[domain.InventoryItemDetails]
