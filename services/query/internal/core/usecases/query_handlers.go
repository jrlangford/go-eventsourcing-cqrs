package usecases

import (
	"fmt"

	"github.com/jrlangford/go-eventsourcing-cqrs/lib/query"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/query/internal/core/secports"
	redis "github.com/redis/go-redis/v9"
)

// An inventoryQueryHandler reads an inventory items list.
type inventoryQueryHandler struct {
	rdb             *redis.Client
	itemListRepo    secports.InventoryItemListRepo
	itemDetailsRepo secports.InventoryItemDetailsRepo
}

// newInventoryQueryHandler constructs a new read model
func newInventoryQueryHandler(
	rdb *redis.Client,
	itemListRepo secports.InventoryItemListRepo,
	itemDetailsRepo secports.InventoryItemDetailsRepo,
) *inventoryQueryHandler {
	return &inventoryQueryHandler{
		rdb:             rdb,
		itemListRepo:    itemListRepo,
		itemDetailsRepo: itemDetailsRepo,
	}
}

// Handle provides a query response corresponding to the provided query message.
func (qh *inventoryQueryHandler) Handle(message query.QueryMessage) (interface{}, error) {

	switch query := message.Query().(type) {

	case *domain.GetInventoryItems:

		items, err := qh.itemListRepo.ReadAll(message.Context())
		if err != nil {
			return nil, err
		}

		itemsList := make([]*secports.InventoryItemListDto, 0, len(items))

		for k, v := range items {
			v.ID = k
			itemsList = append(itemsList, v)
		}

		return itemsList, nil

	case *domain.GetInventoryItemDetails:

		item, err := qh.itemDetailsRepo.Read(message.Context(), query.Uuid)
		if err != nil {
			return nil, err
		}

		item.ID = query.Uuid

		return item, nil

	default:
		return nil, fmt.Errorf("Query can't be processed by this handler: %+v", query)

	}
	return nil, fmt.Errorf("Reached unexpected code section")
}
