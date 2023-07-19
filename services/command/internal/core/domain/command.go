package domain

// A CreateInventoryItem is a self-describing command.
type CreateInventoryItem struct {
	Name string
}

// A DeactivateInventoryItem is a self-describing command.
type DeactivateInventoryItem struct {
}

// RenameInventoryItem renames an inventory item
type RenameInventoryItem struct {
	NewName string
}

// A CheckInItemsToInventory is a self-describing command.
type CheckInItemsToInventory struct {
	Count int
}

// A RemoveItemsFromInventory is a self-describing command.
type RemoveItemsFromInventory struct {
	Count int
}
