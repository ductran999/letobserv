package inventory

import (
	"context"
	"errors"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type Usecase interface {
	InventoryReserve(ctx context.Context, req InventoryReserveInput) (*InventoryReservationView, error)
}
