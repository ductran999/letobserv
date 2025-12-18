package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ductran999/letobserv/internal/domain/repository"
)

type productReduceStockUsecase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productReduceStockUsecase{repo: repo}
}

func (uc *productReduceStockUsecase) ReduceProductStock(ctx context.Context, productID string) error {
	id, err := strconv.Atoi(productID)
	if err != nil || id <= 0 {
		return ErrProductInvalidID
	}

	if err := uc.repo.ReduceStock(ctx, id); err != nil {
		return fmt.Errorf("ReduceProductStock got error: %w", err)
	}

	return nil
}
