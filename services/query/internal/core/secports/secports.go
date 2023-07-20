package secports

import (
	"context"

	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/domain"
)

// ProjectionRepo wraps the methods for the repository of the provided type.
type ProjectionRepo[T any] interface {
	Read(ctx context.Context, dataKey string) (*T, error)
	ReadAll(ctx context.Context) (map[string]*T, error)
}

// InventoryItemListRepo is an alias for an InventoryItemList repo.
type InventoryItemListRepo = ProjectionRepo[domain.InventoryItemList]

// InventoryItemDetailsRepo is an alias for an InventoryItemDetails repo.
type InventoryItemDetailsRepo = ProjectionRepo[domain.InventoryItemDetails]
