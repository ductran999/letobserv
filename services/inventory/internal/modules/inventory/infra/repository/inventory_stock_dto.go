package inventoryrepo

import "time"

type InventoryStockDTO struct {
	ProductID   string    `gorm:"column:product_id;primaryKey;type:uuid" json:"product_id"` // primary key
	TotalQty    int       `gorm:"column:total_qty;not null;check:total_qty >= 0" json:"total_qty"`
	ReservedQty int       `gorm:"column:reserved_qty;not null;default:0;check:reserved_qty >= 0" json:"reserved_qty"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
}

func (InventoryStockDTO) TableName() string { return "inventory_stock" }
