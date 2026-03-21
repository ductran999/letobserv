package inventoryuc

type InventoryReserveInput struct {
	OrderID string
	Items   []ReservationItem
	TTL     int
}

type ReservationItem struct {
	ProductID string
	Quantity  int
}
