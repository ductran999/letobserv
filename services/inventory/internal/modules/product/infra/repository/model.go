package productrepo

import (
	"time"

	product "github.com/ductran999/letobserv/services/inventory/internal/modules/product/domain"
	"github.com/google/uuid"
)

type ProductGORM struct {
	ID          uuid.UUID  `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null;column:name"`
	Description *string    `json:"description" gorm:"type:description;column:description"`
	Price       float64    `json:"price" gorm:"type:numeric(10,2);not null;column:price"`
	CreatedAt   *time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (p *ProductGORM) TableName() string { return "products" }

func (p *ProductGORM) toProductEntity() product.Product {
	return product.Product{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
