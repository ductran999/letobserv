package repository

import (
	"context"

	"github.com/ductran999/letobserv/internal/domain/entity"
)

type ProductRepository interface {
	List(ctx context.Context) ([]entity.Product, error)

	ReduceStock(ctx context.Context, productID int) error
}
