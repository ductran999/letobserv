package order

import "context"

type InventoryReserveRequest struct {
	OrderID string        `json:"order_id"`
	Items   []ReserveItem `json:"items"`
	TTL     int           `json:"ttl_seconds"`
}

type PlaceOrderRequest struct {
	UserID string

	Money           float64
	ShippingAddress string
	ShippingPhone   string

	Items []OrderItem
}

type PlacedOrderOutput struct {
	ID              string
	UserID          string
	OrderNumber     string
	Money           float64
	ShippingAddress string
	ShippingPhone   string
}

type PlacedOrder struct {
	PlacedOrderOutput

	Items []OrderItem
}

type UseCase interface {
	PlacePOrder(ctx context.Context, input PlaceOrderRequest) (*PlacedOrderOutput, error)
	GetOrder(ctx context.Context, id string) (*PlacedOrder, error)
}
