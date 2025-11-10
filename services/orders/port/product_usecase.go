package port

import "context"

type OrderUseCase interface {
	PlacePOrder(ctx context.Context) error
}
