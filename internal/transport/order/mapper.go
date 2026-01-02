package handler

import (
	generated "github.com/ductran999/letobserv/api/generated/orders"
	"github.com/ductran999/letobserv/internal/application/order"
)

// --------- Map OpenAPI struct to Handler --------------

func FromPlaceOrderOpenAPIRequest(body *generated.PlaceOrderJSONRequestBody) *order.PlaceOrderRequest {
	items := make([]order.OrderItem, 0, len(body.Items))
	for _, item := range body.Items {
		items = append(items, order.OrderItem{
			ID:          item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.Id,
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

// --------- Map Handler to OpenAPI struct ------------------

func ToPlacedOrderInfoOpenAPI(output *order.PlacedOrderOutput) *generated.PlacedOrderInfo {
	return &generated.PlacedOrderInfo{
		OrderId:         output.ID,
		UserId:          output.UserID,
		OrderNumber:     output.OrderNumber,
		Money:           output.Money,
		ShippingAddress: output.ShippingAddress,
		ShippingPhone:   output.ShippingPhone,
	}
}
