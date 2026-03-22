package order

import "context"

type InventoryRestful interface {
	Reserve(ctx context.Context, req InventoryReserveRequest) error
	GetProduct(ctx context.Context, id string) error
}
