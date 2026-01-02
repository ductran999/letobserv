package product

import (
	"time"

	"github.com/google/uuid"
)

type ProductDTO struct {
	ID          uuid.UUID  `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null;column:name"`
	Description *string    `json:"description" gorm:"type:description;column:description"`
	Price       float64    `json:"price" gorm:"type:numeric(10,2);not null;column:price"`
	CreatedAt   *time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ProductDTO) TableName() string { return "products" }
