package persistent

import (
	"context"

	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/internal/infrastructure/model"
	"gorm.io/gorm"
)

type productPersistent struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *productPersistent {
	return &productPersistent{db: db}
}

func (repo *productPersistent) ReduceStock(ctx context.Context, productID int) error {
	// Atomic update: only decrease stock if it's > 0
	result := repo.db.WithContext(ctx).Model(&model.Product{}).
		Where("id = ? AND stock > 0", productID).
		Updates(map[string]any{
			"stock":      gorm.Expr("stock - 1"),
			"updated_at": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return usecase.ErrOutOfStock
	}

	return nil
}
