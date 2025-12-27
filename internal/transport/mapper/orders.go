package mapper

import (
	"github.com/ductran999/letobserv/api/generated"
	"github.com/ductran999/letobserv/internal/application/inputs"
	"github.com/ductran999/letobserv/internal/application/outputs"
)

// --------- Map OpenAPI struct to Handler --------------

func FromPlaceOrderOpenAPIRequest(body *generated.PlaceOrderJSONRequestBody) *inputs.PlaceOrderRequest {
	items := make([]inputs.OrderItem, 0, len(body.Items))
	for _, item := range body.Items {
		items = append(items, inputs.OrderItem{
			ID:          item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.Id,
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			UnitPrice:   item.UnitPrice,
		})
	}

	return &inputs.PlaceOrderRequest{
		UserID:          body.UserID.String(),
		Money:           body.Money,
		ShippingAddress: body.ShippingAddress,
		ShippingPhone:   body.ShippingPhone,
		Items:           items,
	}
}

// --------- Map Handler to OpenAPI struct ------------------

func ToPlacedOrderInfoOpenAPI(output *outputs.PlacedOrderOutput) *generated.PlacedOrderInfo {
	return &generated.PlacedOrderInfo{
		OrderId:         output.ID,
		UserId:          output.UserID,
		OrderNumber:     output.OrderNumber,
		Money:           output.Money,
		ShippingAddress: output.ShippingAddress,
		ShippingPhone:   output.ShippingPhone,
	}
}
