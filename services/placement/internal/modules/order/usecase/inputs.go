package orderuc

type OrderItem struct {
	ID          string
	OrderID     string
	ProductID   string
	ProductName string
	Quantity    int
	UnitPrice   float64
}

type PlaceOrderRequest struct {
	UserID string

	Money           float64
	ShippingAddress string
	ShippingPhone   string

	Items []OrderItem
}
