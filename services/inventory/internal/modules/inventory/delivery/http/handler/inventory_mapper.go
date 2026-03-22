package inventoryhttp

import (
	gen "github.com/ductran999/letobserv/services/inventory/api/gen/openapi"
	inventoryuc "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/usecase"
	"github.com/google/uuid"
)

func FromOpenApiInventoryReserveRequest(body *gen.InventoryReserveRequest) *inventoryuc.InventoryReserveInput {
	orderItems := make([]inventoryuc.ReservationItem, 0, len(body.Items))
	for _, item := range body.Items {
		orderItem := inventoryuc.ReservationItem{
			ProductID: item.ProductId.String(),
			Quantity:  item.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	return &inventoryuc.InventoryReserveInput{
		OrderID: body.OrderId.String(),
		Items:   orderItems,
		TTL:     body.TtlSeconds,
	}
}

func ToReservationResponseOpenAPI(reservation *inventoryuc.InventoryReservationView) *gen.ReservationResponse {
	reservationItems := make([]gen.ReservationItem, len(reservation.Items))
	for i, rev := range reservation.Items {
		reservationItems[i] = gen.ReservationItem{
			ProductId: uuid.MustParse(rev.ProductID),
			Quantity:  rev.Quantity,
		}
	}

	return &gen.ReservationResponse{
		ReservationId: reservation.ID,
		OrderId:       reservation.OrderID,
		Items:         reservationItems,
		ExpiredAt:     reservation.TTL,
	}
}
