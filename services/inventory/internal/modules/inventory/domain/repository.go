package inventory

import (
	"context"
)

type InventoryRepo interface {
	GetStock(ctx context.Context, productID string) (*InventoryStock, error)

	IncreaseReserved(ctx context.Context, productID string, quantity int) error

	CreateReservation(ctx context.Context, reservation InventoryReservation) error
}
