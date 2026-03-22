package productsvc

import (
	"context"

	product "github.com/ductran999/letobserv/services/inventory/internal/modules/product/domain"
)

type usecase struct {
	productRepo product.Repository
}

func New(pr product.Repository) product.Usecase {
	if pr == nil {
		panic("product repo is nil")
	}

	return &usecase{
		productRepo: pr,
	}
}

func (uc *usecase) List(ctx context.Context) (*product.ListProductsOutput, error) {
	products, err := uc.productRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	return &product.ListProductsOutput{
		Products: products,
	}, nil
}
