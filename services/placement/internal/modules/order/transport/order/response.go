package orderhttp

import (
	gen "github.com/ductran999/letobserv/services/placement/api/gen/orders"
	order "github.com/ductran999/letobserv/services/placement/internal/modules/order/domain"
)

func toPlacedOrderInfoOpenAPI(output *order.PlacedOrderOutput) *gen.PlacedOrderInfo {
	return &gen.PlacedOrderInfo{
		OrderId:         output.ID,
		UserId:          output.UserID,
		OrderNumber:     output.OrderNumber,
		Money:           output.Money,
		ShippingAddress: output.ShippingAddress,
		ShippingPhone:   output.ShippingPhone,
	}
}

func toPlacedOrderDetailsOpenAPI(o *order.PlacedOrder) gen.PlacedOrderDetails {
	orderItems := make([]any, len(o.Items))
	for i, item := range o.Items {
		orderItems[i] = item
	}

	return gen.PlacedOrderDetails{
		OrderId:         o.ID,
		Money:           o.Money,
		OrderItems:      orderItems,
		OrderNumber:     o.OrderNumber,
		ShippingAddress: o.ShippingAddress,
		ShippingPhone:   o.ShippingPhone,
		UserId:          o.UserID,
	}
}
