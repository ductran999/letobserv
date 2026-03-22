package orderrepo

import "time"

type OrderDTO struct {
	ID          string `gorm:"column:id;type:uuid;primaryKey"`
	OrderNumber string `gorm:"column:order_number;type:varchar(50);uniqueIndex;not null"`
	UserID      string `gorm:"column:user_id;type:uuid;index;not null"`

	Status string  `gorm:"column:status;type:varchar(30);index;not null"`
	Money  float64 `gorm:"column:total_amount;not null"`

	ShippingAddress      string `gorm:"column:shipping_address;type:text"`
	ShippingPhone        string `gorm:"column:shipping_phone;type:varchar(20)"`
	PaymentTransactionID string `gorm:"column:payment_transaction_id;type:varchar(100);index"`
	ShippingTrackingID   string `gorm:"column:shipping_tracking_id;type:varchar(100);index"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (o *OrderDTO) TableName() string { return "orders" }

type OrderItemDTO struct {
	ID          string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrderID     string  `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID   string  `gorm:"type:uuid;not null" json:"product_id"`
	ProductName string  `gorm:"type:varchar(255);not null" json:"product_name"`
	Quantity    int     `gorm:"not null" json:"quantity"`
	UnitPrice   float64 `gorm:"not null" json:"unit_price"`
}

func (o *OrderItemDTO) TableName() string { return "order_items" }
