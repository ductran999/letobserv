package inventory

import "time"

type InventoryReservationView struct {
	ID      string
	OrderID string
	Items   []ReservationItem
	TTL     time.Time
}
