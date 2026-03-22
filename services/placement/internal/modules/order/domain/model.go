package order

import "time"

type OrderStatus string

const (
	OrderPendingStatus OrderStatus = "pending"
)

type OrderItem struct {
	ID          string
	OrderID     string
	ProductID   string
	ProductName string
	Quantity    int
	UnitPrice   float64
}

type Order struct {
	// Identity
	ID          string
	OrderNumber string
	UserID      string

	Status               OrderStatus
	Money                float64
	ShippingAddress      string
	ShippingPhone        string
	PaymentTransactionID string
	ShippingTrackingID   string

	Items []OrderItem

	CreateAt time.Time
	UpdateAt time.Time
}
