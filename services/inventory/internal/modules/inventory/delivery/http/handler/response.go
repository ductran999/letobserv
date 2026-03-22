package inventoryhttp

import (
	gen "github.com/ductran999/letobserv/services/inventory/api/gen/openapi"
	inventory "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/domain"
	"github.com/google/uuid"
)

func toReservationResponseOpenAPI(reservation *inventory.InventoryReservationView) *gen.ReservationResponse {
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
