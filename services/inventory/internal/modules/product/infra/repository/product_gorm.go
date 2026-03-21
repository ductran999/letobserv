package productrepo

import (
	"context"

	product "github.com/ductran999/letobserv/services/inventory/internal/modules/product/domain"
	"gorm.io/gorm"
)

type productPersistent struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *productPersistent {
	return &productPersistent{db: db}
}

func (repo *productPersistent) List(ctx context.Context) ([]product.Product, error) {
	queryResult := make([]ProductGORM, 0)
	if err := repo.db.WithContext(ctx).Table((&ProductGORM{}).TableName()).Find(&queryResult).Error; err != nil {
		return nil, err
	}

	product := make([]product.Product, len(queryResult))
	for i, r := range queryResult {
		product[i] = r.ToProductEntity()
	}

	return product, nil
}

func (repo *productPersistent) GetByID(ctx context.Context, id string) (*product.Product, error) {
	queryResult := ProductGORM{}
	if err := repo.db.WithContext(ctx).
		Table((&ProductGORM{}).TableName()).
		Find(&queryResult).Error; err != nil {
		return nil, err
	}

	p := queryResult.ToProductEntity()

	return &p, nil
}
