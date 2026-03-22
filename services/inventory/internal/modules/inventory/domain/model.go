package inventory

import "time"

type ReservationStatus string

const (
	ReservationStatusReserved  ReservationStatus = "RESERVED"
	ReservationStatusConfirmed ReservationStatus = "CONFIRMED"
	ReservationStatusReleased  ReservationStatus = "RELEASED"
)

type InventoryReservation struct {
	ID string

	OrderID   string
	ProductID string

	Quantity int
	Status   ReservationStatus

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
}

type InventoryStock struct {
	ProductID   string
	TotalQty    int
	ReservedQty int
	UpdatedAt   time.Time
}
