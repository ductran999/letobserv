package entity

import "time"

type Order struct {
	// Identity
	ID          string
	OrderNumber string
	UserID      string

	Status               string
	Money                float64
	ShippingAddress      string
	ShippingPhone        string
	PaymentTransactionID string
	ShippingTrackingID   string

	Items []OrderItem

	CreateAt time.Time
	UpdateAt time.Time
}
