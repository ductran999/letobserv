package order

import (
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
}
