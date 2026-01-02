package inventory

import "time"

type ReservationStatus string

const (
	ReservationReserved  ReservationStatus = "RESERVED"
	ReservationConfirmed ReservationStatus = "CONFIRMED"
	ReservationReleased  ReservationStatus = "RELEASED"
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
