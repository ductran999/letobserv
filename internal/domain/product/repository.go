package product

import (
	"context"
)

type ProductRepository interface {
	List(ctx context.Context) ([]Product, error)
}
