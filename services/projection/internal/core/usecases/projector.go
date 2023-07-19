package usecases

import (
	"fmt"

	es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/core/domain"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/projection/internal/core/secports"
)

// A Projector is a façade to the views.
type Projector struct {
	invItemListRepo    secports.InventoryItemListRepo
	invItemDetailsRepo secports.InventoryItemDetailsRepo
	eventBus           *es.InternalEventBus
}

// NewProjector is self-describing.
func NewProjector(
	invItemListRepo secports.InventoryItemListRepo,
	invItemDetailsRepo secports.InventoryItemDetailsRepo,
) *Projector {

	// Create views
	listView := newInventoryListView(invItemListRepo)
	detailView := newInventoryItemDetailView(invItemDetailsRepo)

	eventBus := es.NewInternalEventBus()

	eventBus.AddHandler(listView,
		&domain.InventoryItemCreated{},
		&domain.InventoryItemRenamed{},
		&domain.InventoryItemDeactivated{},
	)

	eventBus.AddHandler(detailView,
		&domain.InventoryItemCreated{},
		&domain.InventoryItemRenamed{},
		&domain.InventoryItemDeactivated{},
		&domain.ItemsRemovedFromInventory{},
		&domain.ItemsCheckedIntoInventory{},
	)

	return &Projector{
		invItemListRepo:    invItemListRepo,
		invItemDetailsRepo: invItemDetailsRepo,
		eventBus:           eventBus,
	}
}

// ProjectEvent projects the incoming events using their corresponding views, if present.
func (p *Projector) ProjectEvent(
	aggregateID string,
	domainEvent interface{},
	version *int,
	metadata []byte,
) error {
	fmt.Printf("Domain event: %+v\n", domainEvent)

	em := es.NewEventMessage(aggregateID, domainEvent, version)
	em.SetUserMetadata(metadata)
	p.eventBus.PublishEvent(em)
	// TODO: Handle errors
	return nil
}
