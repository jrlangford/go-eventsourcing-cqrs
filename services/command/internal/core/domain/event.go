package domain

// An InventoryItemCreated is a self-dscribing event.
type InventoryItemCreated struct {
	Name string
}

// An InventoryItemRenamed is a self-dscribing event.
type InventoryItemRenamed struct {
	NewName string
}

// An InventoryItemDeactivated is a self-dscribing event.
type InventoryItemDeactivated struct {
}

// An ItemsRemovedFromInventory is a self-dscribing event.
type ItemsRemovedFromInventory struct {
	Count int
}

// An ItemsCheckedIntoInventory is a self-dscribing event.
type ItemsCheckedIntoInventory struct {
	Count int
}
