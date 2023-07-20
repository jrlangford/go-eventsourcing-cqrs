package domain

// An InventoryItemDetails is a self-describing read model.
type InventoryItemDetails struct {
	ID           string `msgpack:"-"`
	Name         string
	CurrentCount int
	Version      int
}

// An InventoryItemList is a self-describing read model.
type InventoryItemList struct {
	ID   string `msgpack:"-"`
	Name string
}
