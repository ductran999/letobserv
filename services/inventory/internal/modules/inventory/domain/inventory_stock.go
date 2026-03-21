package inventory

import "time"

type InventoryStock struct {
	ProductID   string
	TotalQty    int
	ReservedQty int
	UpdatedAt   time.Time
}
