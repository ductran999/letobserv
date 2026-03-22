package inventoryuc

import (
	"context"
	"time"

	"github.com/ductran999/letobserv/pkg/txmanager"
	inventory "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/domain"
	product "github.com/ductran999/letobserv/services/inventory/internal/modules/product/domain"
	"github.com/google/uuid"
)

type inventoryUC struct {
	txManger txmanager.TransactionManager

	productRepo   product.Repository
	inventoryRepo inventory.InventoryRepo
}

func NewInventoryUsecase(
	txManger txmanager.TransactionManager,
	productRepo product.Repository,
	inventoryRepo inventory.InventoryRepo,
) inventory.Usecase {
	return &inventoryUC{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
		txManger:      txManger,
	}
}

func (uc *inventoryUC) InventoryReserve(
	ctx context.Context,
	req inventory.InventoryReserveInput,
) (*inventory.InventoryReservationView, error) {
	reservations := make([]inventory.InventoryReservation, 0, len(req.Items))

	err := uc.txManger.Do(ctx, func(txContext context.Context) error {
		now := time.Now()
		reservationID := uuid.New().String()
		for _, product := range req.Items {
			if err := uc.inventoryRepo.IncreaseReserved(txContext, product.ProductID, product.Quantity); err != nil {
				return err
			}

			reservation := inventory.InventoryReservation{
				ID:        reservationID,
				OrderID:   req.OrderID,
				ProductID: product.ProductID,
				Quantity:  product.Quantity,
				Status:    inventory.ReservationStatusReserved,
				CreatedAt: now,
				UpdatedAt: now,
				ExpiredAt: now.Add(time.Duration(req.TTL)),
			}
			if err := uc.inventoryRepo.CreateReservation(txContext, reservation); err != nil {
				return err
			}
			reservations = append(reservations, reservation)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return toInventoryReservationView(reservations), nil
}
