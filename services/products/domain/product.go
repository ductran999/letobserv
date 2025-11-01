package domain

import (
	"errors"
	"time"
)

var (
	ErrProductInvalidID = errors.New("invalid product id")
	ErrOutOfStock       = errors.New("out of stock")
)

type Product struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null;column:name"`
	SKU       string    `json:"sku" gorm:"type:varchar(100);uniqueIndex;not null;column:sku"`
	Price     float64   `json:"price" gorm:"type:numeric(10,2);not null;column:price"`
	Stock     int       `json:"stock" gorm:"not null;default:0;column:stock"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
}

func (Product) TableName() string {
	return "products"
}
