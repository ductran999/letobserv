package inventoryuc

import (
	"time"
)

// type ListProductsOutput struct {
// 	Products []product.Product
// }

type InventoryReservationView struct {
	ID      string
	OrderID string
	Items   []ReservationItem
	TTL     time.Time
}
