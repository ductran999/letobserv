package inventoryuc

import inventory "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/domain"

func toInventoryReservationView(reservations []inventory.InventoryReservation) *inventory.InventoryReservationView {
	reserveItems := make([]inventory.ReservationItem, 0)
	for _, rev := range reservations {
		reserveItems = append(reserveItems, inventory.ReservationItem{
			ProductID: rev.ProductID,
			Quantity:  rev.Quantity,
		})
	}

	return &inventory.InventoryReservationView{
		ID:      reservations[0].ID,
		OrderID: reservations[0].OrderID,
		Items:   reserveItems,
		TTL:     reservations[0].ExpiredAt,
	}
}
