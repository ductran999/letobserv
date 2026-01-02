package inventory

import (
	"time"

	"github.com/ductran999/letobserv/internal/domain/product"
)

type ListProductsOutput struct {
	Products []product.Product
}

type InventoryReservationView struct {
	ID      string
	OrderID string
	Items   []ReservationItem
	TTL     time.Time
}
