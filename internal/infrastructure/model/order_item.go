package model

type OrderItemDTO struct {
	ID          string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrderID     string  `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID   string  `gorm:"type:uuid;not null" json:"product_id"`
	ProductName string  `gorm:"type:varchar(255);not null" json:"product_name"`
	Quantity    int     `gorm:"not null" json:"quantity"`
	UnitPrice   float64 `gorm:"not null" json:"unit_price"`
}

func (o *OrderItemDTO) TableName() string { return "order_items" }
