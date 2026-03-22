package orderhttp

import (
	gen "github.com/ductran999/letobserv/services/placement/api/gen/orders"
	order "github.com/ductran999/letobserv/services/placement/internal/modules/order/domain"
)

func toPlaceOrderRequest(body *gen.PlaceOrderJSONRequestBody) *order.PlaceOrderRequest {
	items := make([]order.OrderItem, 0, len(body.Items))
	for _, item := range body.Items {
		items = append(items, order.OrderItem{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			UnitPrice:   item.UnitPrice,
		})
	}

	return &order.PlaceOrderRequest{
		UserID:          body.UserID.String(),
		Money:           body.Money,
		ShippingAddress: body.ShippingAddress,
		ShippingPhone:   body.ShippingPhone,
		Items:           items,
	}
}
