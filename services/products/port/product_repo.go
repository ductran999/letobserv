package port

import (
	"context"
)

type ProductRepo interface {
	ReduceStock(ctx context.Context, productID int) error
}
