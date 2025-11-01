package repo

import (
	"context"

	"github.com/ductran999/letobserv/services/products/domain"
	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *productRepo {
	return &productRepo{db: db}
}

func (repo *productRepo) ReduceStock(ctx context.Context, productID int) error {
	// Atomic update: only decrease stock if it's > 0
	result := repo.db.WithContext(ctx).Model(&domain.Product{}).
		Where("id = ? AND stock > 0", productID).
		Updates(map[string]any{
			"stock":      gorm.Expr("stock - 1"),
			"updated_at": gorm.Expr("NOW()"),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrOutOfStock
	}

	return nil
}
