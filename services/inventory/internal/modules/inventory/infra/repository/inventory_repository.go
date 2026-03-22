package inventoryrepo

import (
	"context"

	"github.com/ductran999/letobserv/pkg/txmanager"
	inventory "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/domain"
	"gorm.io/gorm"
)

type inventoryPersistent struct {
	db *gorm.DB
}

func NewInventoryRepo(db *gorm.DB) inventory.InventoryRepo {
	return &inventoryPersistent{db: db}
}

func (r *inventoryPersistent) GetStock(ctx context.Context, productID string) (*inventory.InventoryStock, error) {
	txDB, err := txmanager.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	queryResult := InventoryStockDTO{}
	if err := txDB.WithContext(ctx).Table((&InventoryStockDTO{}).TableName()).Find(&queryResult).Error; err != nil {
		return nil, err
	}

	stock := inventory.InventoryStock{
		ProductID:   queryResult.ProductID,
		TotalQty:    queryResult.TotalQty,
		ReservedQty: queryResult.ReservedQty,
		UpdatedAt:   queryResult.UpdatedAt,
	}

	return &stock, nil
}

func (r *inventoryPersistent) IncreaseReserved(ctx context.Context, productID string, quantity int) error {
	txDB, err := txmanager.GetTx(ctx)
	if err != nil {
		return err
	}

	result := txDB.WithContext(ctx).
		Model(&InventoryStockDTO{}).
		Where(
			"product_id = ? AND total_qty - reserved_qty >= ?",
			productID,
			quantity,
		).
		UpdateColumn(
			"reserved_qty",
			gorm.Expr("reserved_qty + ?", quantity),
		)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return inventory.ErrInsufficientStock
	}

	return nil
}

func (r *inventoryPersistent) CreateReservation(ctx context.Context, reservation inventory.InventoryReservation) error {
	txDB, err := txmanager.GetTx(ctx)
	if err != nil {
		return err
	}

	dto := InventoryReservationDTO{
		OrderID:   reservation.ID,
		ProductID: reservation.ProductID,
		Quantity:  reservation.Quantity,
		Status:    inventory.ReservationStatusReserved,
		CreatedAt: reservation.CreatedAt,
		UpdatedAt: reservation.UpdatedAt,
		ExpiredAt: reservation.ExpiredAt,
	}

	if err := txDB.WithContext(ctx).
		Create(&dto).Error; err != nil {
		return err
	}

	return nil
}
