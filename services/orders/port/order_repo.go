package port

import (
	"context"
)

type OrderRepo interface {
	Create(ctx context.Context, productID int) error
}
