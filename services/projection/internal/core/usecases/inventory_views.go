package usecases

import (
	"context"

	es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/core/secports"
)

// An inventoryListView projects an inventory items list.
type inventoryListView struct {
	repo secports.InventoryItemListRepo
	ctx  context.Context
}

// newInventoryListView is self-describing.
func newInventoryListView(repo secports.InventoryItemListRepo) *inventoryListView {
	return &inventoryListView{
		repo: repo,
		ctx:  context.Background(),
	}
}

// Handle projects the model from the incoming events.
func (v *inventoryListView) Handle(message es.EventMessage) {

	switch event := message.Event().(type) {

	case *domain.InventoryItemCreated:

		v.repo.Write(
			v.ctx,
			message.AggregateID(),
			&secports.InventoryItemListDto{Name: event.Name},
		)

	case *domain.InventoryItemRenamed:

		v.repo.Update(
			v.ctx,
			message.AggregateID(),
			func(dto *secports.InventoryItemListDto) {
				dto.Name = event.NewName
			},
		)

	case *domain.InventoryItemDeactivated:

		v.repo.Del(v.ctx, message.AggregateID())

	}
}

// inventoryItemDetailView projects invntory item details.
type inventoryItemDetailView struct {
	repo secports.InventoryItemDetailsRepo
	ctx  context.Context
}

// newInventoryListView is self-describing.
func newInventoryItemDetailView(repo secports.InventoryItemDetailsRepo) *inventoryItemDetailView {
	return &inventoryItemDetailView{
		repo: repo,
		ctx:  context.Background(),
	}
}

// Handle projects the model from the incoming events.
func (v *inventoryItemDetailView) Handle(message es.EventMessage) {

	switch event := message.Event().(type) {

	case *domain.InventoryItemCreated:

		v.repo.Write(
			v.ctx,
			message.AggregateID(),
			&secports.InventoryItemDetailsDto{
				Name:    event.Name,
				Version: 0,
			},
		)

	case *domain.InventoryItemRenamed:

		v.repo.Update(
			v.ctx,
			message.AggregateID(),
			func(dto *secports.InventoryItemDetailsDto) {
				dto.Name = event.NewName
				dto.Version = *message.Version()
			},
		)

	case *domain.ItemsRemovedFromInventory:

		v.repo.Update(
			v.ctx,
			message.AggregateID(),
			func(dto *secports.InventoryItemDetailsDto) {
				dto.CurrentCount -= event.Count
				dto.Version = *message.Version()
			},
		)

	case *domain.ItemsCheckedIntoInventory:

		v.repo.Update(
			v.ctx,
			message.AggregateID(),
			func(dto *secports.InventoryItemDetailsDto) {
				dto.CurrentCount += event.Count
				dto.Version = *message.Version()
			},
		)

	case *domain.InventoryItemDeactivated:

		v.repo.Del(v.ctx, message.AggregateID())

	}
}
