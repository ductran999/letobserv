package product

import (
	"context"
)

type ProductRepository interface {
	List(ctx context.Context) ([]Product, error)

	GetByID(ctx context.Context, id string) (*Product, error)
}
