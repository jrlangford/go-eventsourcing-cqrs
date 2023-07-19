package secadapters

import (
	"reflect"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"
	"github.com/jrlangford/go-eventsourcing-cqrs/lib/jsontools"
	"github.com/jrlangford/go-eventsourcing-cqrs/services/command/internal/core/domain"
)

// An InventoryItemRepo reads and writes InventoryItem-related events to the event store.
type InventoryItemRepo struct {
	baseRepo *es.AggregateRepoBase[domain.InventoryItem, *domain.InventoryItem]
}

// NewInventoryItemRepo is sel-describing.
func NewInventoryItemRepo(eventStore *esdb.Client) (*InventoryItemRepo, error) {

	r, err := es.NewAggregateRepoBase[domain.InventoryItem, *domain.InventoryItem](eventStore)
	if err != nil {
		return nil, err
	}

	repo := &InventoryItemRepo{
		baseRepo: r,
	}

	aggregateFactory := es.NewDelegateAggregateFactory[domain.InventoryItem, *domain.InventoryItem](
		func(id string) *domain.InventoryItem { return domain.NewInventoryItem(id) },
	)
	repo.baseRepo.SetAggregateFactory(aggregateFactory)

	streamNameDelegate := es.NewAggregateStreamNamer(reflect.TypeOf(((*domain.InventoryItem)(nil))).Elem().Name())
	repo.baseRepo.SetStreamNameDelegate(streamNameDelegate)

	eventFactory := es.NewDelegateEventFactory()
	eventFactory.RegisterDelegate(&domain.InventoryItemCreated{}, jsontools.JsonToTypeInterface[domain.InventoryItemCreated])
	eventFactory.RegisterDelegate(&domain.InventoryItemRenamed{}, jsontools.JsonToTypeInterface[domain.InventoryItemRenamed])
	eventFactory.RegisterDelegate(&domain.InventoryItemDeactivated{}, jsontools.JsonToTypeInterface[domain.InventoryItemDeactivated])
	eventFactory.RegisterDelegate(&domain.ItemsRemovedFromInventory{}, jsontools.JsonToTypeInterface[domain.ItemsRemovedFromInventory])
	eventFactory.RegisterDelegate(&domain.ItemsCheckedIntoInventory{}, jsontools.JsonToTypeInterface[domain.ItemsCheckedIntoInventory])
	repo.baseRepo.SetEventFactory(eventFactory)

	return repo, nil
}

// Load builds an inventory item by loading its events.
func (r *InventoryItemRepo) Load(id string) (*domain.InventoryItem, error) {
	ar, err := r.baseRepo.Load(id)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

// Save persists an inventory item by saving its events.
func (r *InventoryItemRepo) Save(aggregate es.AggregateRoot, expectedVersion *int) error {
	return r.baseRepo.Save(aggregate, expectedVersion)
}
