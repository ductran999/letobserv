package product

import (
	"context"

	"github.com/ductran999/letobserv/internal/domain/product"
	"gorm.io/gorm"
)

type productPersistent struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *productPersistent {
	return &productPersistent{db: db}
}

func (repo *productPersistent) List(ctx context.Context) ([]product.Product, error) {
	queryResult := make([]ProductDTO, 0)
	if err := repo.db.WithContext(ctx).Table((&ProductDTO{}).TableName()).Find(&queryResult).Error; err != nil {
		return nil, err
	}

	product := make([]product.Product, len(queryResult))
	for i, r := range queryResult {
		product[i] = repo.toProductEntity(&r)
	}

	return product, nil
}

func (repo *productPersistent) GetByID(ctx context.Context, id string) (*product.Product, error) {
	queryResult := ProductDTO{}
	if err := repo.db.WithContext(ctx).
		Table((&ProductDTO{}).TableName()).
		Find(&queryResult).Error; err != nil {
		return nil, err
	}

	p := repo.toProductEntity(&queryResult)

	return &p, nil
}

func (repo *productPersistent) toProductEntity(p *ProductDTO) product.Product {
	return product.Product{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
