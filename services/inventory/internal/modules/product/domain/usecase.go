package product

import (
	"context"
)

type ListProductsOutput struct {
	Products []Product
}

type Usecase interface {
	List(ctx context.Context) (*ListProductsOutput, error)
}
