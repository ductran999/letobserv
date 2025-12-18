package repository

import "context"

type OrderRepository interface {
	Create(ctx context.Context, productID int) error
}
