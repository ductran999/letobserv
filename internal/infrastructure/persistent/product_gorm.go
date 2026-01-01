package persistent

import (
	"context"

	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/internal/domain/entity"
	"github.com/ductran999/letobserv/internal/infrastructure/model"
	"gorm.io/gorm"
)

type productPersistent struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *productPersistent {
	return &productPersistent{db: db}
}

func (repo *productPersistent) List(ctx context.Context) ([]entity.Product, error) {
	queryResult := make([]model.Product, 0)
	if err := repo.db.WithContext(ctx).Table((&model.Product{}).TableName()).Find(&queryResult).Error; err != nil {
		return nil, err
	}

	products := make([]entity.Product, len(queryResult))
	for i, r := range queryResult {
		products[i] = repo.toProductEntity(&r)
	}

	return products, nil
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

func (repo *productPersistent) toProductEntity(p *model.Product) entity.Product {
	return entity.Product{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
