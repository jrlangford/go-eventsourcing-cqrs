package secports

import (
	"context"
)

// ProjectionRepo wraps the methods for the repository of the provided type.
type ProjectionRepo[T any] interface {
	Read(ctx context.Context, dataKey string) (*T, error)
	ReadAll(ctx context.Context) (map[string]*T, error)
}

// InventoryItemListRepo is an alias for an InventoryItemListDto repo.
type InventoryItemListRepo = ProjectionRepo[InventoryItemListDto]

// InventoryItemDetailsRepo is an alias for an InventoryItemDetailsDto repo.
type InventoryItemDetailsRepo = ProjectionRepo[InventoryItemDetailsDto]
