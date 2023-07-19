package secports

// An InventoryItemDetailsDto is a read model.
type InventoryItemDetailsDto struct {
	ID           string `msgpack:"-"`
	Name         string
	CurrentCount int
	Version      int
}

// An InventoryItemListDto is a read model.
type InventoryItemListDto struct {
	ID   string `msgpack:"-"`
	Name string
}
