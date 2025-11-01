package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ductran999/letobserv/services/products/domain"
	"github.com/ductran999/letobserv/services/products/port"
)

type productUseCase struct {
	repo port.ProductRepo
}

func NewProductUseCase(repo port.ProductRepo) *productUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) ReduceProductStock(ctx context.Context, productID string) error {
	id, err := strconv.Atoi(productID)
	if err != nil || id <= 0 {
		return domain.ErrProductInvalidID
	}

	if err := uc.repo.ReduceStock(ctx, id); err != nil {
		return fmt.Errorf("ReduceProductStock got error: %w", err)
	}

	return nil
}
