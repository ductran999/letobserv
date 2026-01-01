package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ductran999/letobserv/internal/application/errs"
	"github.com/ductran999/letobserv/internal/application/outputs"
	"github.com/ductran999/letobserv/internal/domain/repository"
)

type InventoryUsecase interface {
	ListProducts(ctx context.Context) (*outputs.ListProductsOutput, error)

	ReduceProductStock(ctx context.Context, productID string) error
}

type inventoryUC struct {
	productRepo repository.ProductRepository
}

func NewInventoryUsecase(productRepo repository.ProductRepository) InventoryUsecase {
	return &inventoryUC{productRepo: productRepo}
}

func (uc *inventoryUC) ListProducts(ctx context.Context) (*outputs.ListProductsOutput, error) {
	products, err := uc.productRepo.List(ctx)
	if err != nil {
		return nil, errs.Internal(err)
	}

	return &outputs.ListProductsOutput{
		Products: products,
	}, nil
}

func (uc *inventoryUC) ReduceProductStock(ctx context.Context, productID string) error {
	id, err := strconv.Atoi(productID)
	if err != nil || id <= 0 {
		return ErrProductInvalidID
	}

	if err := uc.productRepo.ReduceStock(ctx, id); err != nil {
		return fmt.Errorf("ReduceProductStock got error: %w", err)
	}

	return nil
}
