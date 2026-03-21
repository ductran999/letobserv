package handler

import (
	gen "github.com/ductran999/letobserv/services/placement/api/gen/orders"
	orderuc "github.com/ductran999/letobserv/services/placement/internal/modules/order/usecase"
)

// --------- Map OpenAPI struct to Handler --------------

func FromPlaceOrderOpenAPIRequest(body *gen.PlaceOrderJSONRequestBody) *orderuc.PlaceOrderRequest {
	items := make([]orderuc.OrderItem, 0, len(body.Items))
	for _, item := range body.Items {
		items = append(items, orderuc.OrderItem{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			UnitPrice:   item.UnitPrice,
		})
	}

	return &orderuc.PlaceOrderRequest{
		UserID:          body.UserID.String(),
		Money:           body.Money,
		ShippingAddress: body.ShippingAddress,
		ShippingPhone:   body.ShippingPhone,
		Items:           items,
	}
}

// --------- Map Handler to OpenAPI struct ------------------

func ToPlacedOrderInfoOpenAPI(output *orderuc.PlacedOrderOutput) *gen.PlacedOrderInfo {
	return &gen.PlacedOrderInfo{
		OrderId:         output.ID,
		UserId:          output.UserID,
		OrderNumber:     output.OrderNumber,
		Money:           output.Money,
		ShippingAddress: output.ShippingAddress,
		ShippingPhone:   output.ShippingPhone,
	}
}
