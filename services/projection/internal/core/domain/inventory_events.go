package domain

// An InventoryItemCreated is a self-describing event.
type InventoryItemCreated struct {
	Name string
}

// An InventoryItemRenamed is a self-describing event.
type InventoryItemRenamed struct {
	NewName string
}

// An InventoryItemDeactivated is a self-describing event.
type InventoryItemDeactivated struct {
}

// An ItemsRemovedFromInventory is a self-describing event.
type ItemsRemovedFromInventory struct {
	Count int
}

// An ItemsCheckedIntoInventory is a self-describing event.
type ItemsCheckedIntoInventory struct {
	Count int
}
