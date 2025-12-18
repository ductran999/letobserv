package repository

import "context"

type ProductRepository interface {
	ReduceStock(ctx context.Context, productID int) error
}
