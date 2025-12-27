package inputs

type PlaceOrderRequest struct {
	UserID string

	Money           float64
	ShippingAddress string
	ShippingPhone   string

	Items []OrderItem
}
