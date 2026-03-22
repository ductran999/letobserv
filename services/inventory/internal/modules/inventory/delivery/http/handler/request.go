package inventoryhttp

import (
	gen "github.com/ductran999/letobserv/services/inventory/api/gen/openapi"
	inventory "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/domain"
)

func FromOpenApiInventoryReserveRequest(body *gen.InventoryReserveRequest) *inventory.InventoryReserveInput {
	orderItems := make([]inventory.ReservationItem, 0, len(body.Items))
	for _, item := range body.Items {
		orderItem := inventory.ReservationItem{
			ProductID: item.ProductId.String(),
			Quantity:  item.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	return &inventory.InventoryReserveInput{
		OrderID: body.OrderId.String(),
		Items:   orderItems,
		TTL:     body.TtlSeconds,
	}
}
