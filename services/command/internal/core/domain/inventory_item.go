package domain

import (
	"errors"

	es "github.com/jrlangford/go-eventsourcing-cqrs/lib/eventsourcing"
)

// An InventoryItem is a self-describing aggregate.
type InventoryItem struct {
	*es.AggregateBase
	activated bool
	count     int
}

// NewInventoryItem constructs a new inventory item aggregate.
func NewInventoryItem(id string) *InventoryItem {
	i := &InventoryItem{
		AggregateBase: es.NewAggregateBase(id),
	}

	return i
}

// Create cretes an inventory item.
func (a *InventoryItem) Create(name string) error {
	if name == "" {
		return errors.New("the name can not be empty")
	}

	a.Apply(es.NewEventMessage(a.AggregateID(),
		&InventoryItemCreated{Name: name},
		es.Int(a.CurrentVersion())), true)

	return nil
}

// ChangeName changes the name of the item.
func (a *InventoryItem) ChangeName(newName string) error {
	if newName == "" {
		return errors.New("the name can not be empty")
	}

	a.Apply(es.NewEventMessage(a.AggregateID(),
		&InventoryItemRenamed{NewName: newName},
		es.Int(a.CurrentVersion())), true)

	return nil
}

// Remove removes items from inventory.
func (a *InventoryItem) Remove(count int) error {
	if count <= 0 {
		return errors.New("can't remove negative count from inventory")
	}

	if a.count-count < 0 {
		return errors.New("can't remove more items from inventory than the number of items in inventory")
	}

	a.Apply(es.NewEventMessage(a.AggregateID(),
		&ItemsRemovedFromInventory{Count: count},
		es.Int(a.CurrentVersion())), true)

	return nil
}

// CheckIn adds items to inventory.
func (a *InventoryItem) CheckIn(count int) error {
	if count <= 0 {
		return errors.New("must have a count greater than 0 to add to inventory")
	}

	a.Apply(es.NewEventMessage(a.AggregateID(),
		&ItemsCheckedIntoInventory{Count: count},
		es.Int(a.CurrentVersion())), true)

	return nil
}

// Deactivate deactivates the inventory item.
func (a *InventoryItem) Deactivate() error {
	if !a.activated {
		return errors.New("already deactivated")
	}

	a.Apply(es.NewEventMessage(a.AggregateID(),
		&InventoryItemDeactivated{},
		es.Int(a.CurrentVersion())), true)

	return nil
}

// Apply handles the logic of events on the aggregate.
func (a *InventoryItem) Apply(message es.EventMessage, isNew bool) {

	if isNew {
		a.TrackChange(message)
	}

	switch ev := message.Event().(type) {

	case *InventoryItemCreated:
		a.activated = true

	case *InventoryItemDeactivated:
		a.activated = false

	case *ItemsRemovedFromInventory:
		a.count -= ev.Count

	case *ItemsCheckedIntoInventory:
		a.count += ev.Count

	}

}
