package inventory

import (
	"time"

	"github.com/ductran999/letobserv/internal/domain/inventory"
)

type InventoryReservationDTO struct {
	ID string `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`

	OrderID   string `gorm:"column:order_id;type:uuid;not null;uniqueIndex:ux_order_product"`
	ProductID string `gorm:"column:product_id;type:uuid;not null;uniqueIndex:ux_order_product"`

	Quantity int `gorm:"column:quantity;not null;check:quantity > 0"`

	Status inventory.ReservationStatus `gorm:"column:status;type:text"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
	ExpiredAt time.Time `gorm:"column:expired_at;not null"`
}

func (InventoryReservationDTO) TableName() string { return "inventory_reservations" }
