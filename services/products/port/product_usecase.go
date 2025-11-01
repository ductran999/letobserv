package port

import "context"

type ProductUseCase interface {
	ReduceProductStock(ctx context.Context, productID string) error
}
