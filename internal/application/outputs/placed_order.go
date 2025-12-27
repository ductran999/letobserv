package outputs

type PlacedOrderOutput struct {
	ID              string
	UserID          string
	OrderNumber     string
	Money           float64
	ShippingAddress string
	ShippingPhone   string
}
