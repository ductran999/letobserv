package service

import (
	"context"
)

type InventoryReserveRequest struct {
	OrderID string        `json:"order_id"`
	Items   []ReserveItem `json:"items"`
	TTL     int           `json:"ttl_seconds"`
}

type ReserveItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type InventoryRestful interface {
	Reserve(ctx context.Context, req InventoryReserveRequest) error
}
