package secports

import (
	"context"
)

// ProjectionRepo wraps the methods for the repository of the provided type.
type ProjectionRepo[T any] interface {
	Write(ctx context.Context, dataKey string, val *T) error
	Update(ctx context.Context, dataKey string, updater func(*T)) error
	Del(ctx context.Context, dataKey string) error
}

// InventoryItemListRepo is an alias for an InventoryItemListDto projection repo.
type InventoryItemListRepo = ProjectionRepo[InventoryItemListDto]

// InventoryItemDetailsRepo is an alias for an InventoryItemDetailsDto projection repo.
type InventoryItemDetailsRepo = ProjectionRepo[InventoryItemDetailsDto]
