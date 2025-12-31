package persistent

import (
	"context"

	"github.com/ductran999/letobserv/internal/domain/entity"
	"github.com/ductran999/letobserv/internal/domain/repository"
	"github.com/ductran999/letobserv/internal/infrastructure/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderPersistent struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) repository.OrderRepository {
	return &orderPersistent{db: db}
}

func (r *orderPersistent) Create(ctx context.Context, order *entity.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		orderDTO := r.mapToOrderDTO(order)
		if err := tx.WithContext(ctx).Table(orderDTO.TableName()).Create(orderDTO).Error; err != nil {
			return err
		}

		itemDTOs := r.mapToOrderItemDTOs(order.Items, orderDTO.ID)
		if err := tx.WithContext(ctx).Table((&model.OrderItemDTO{}).TableName()).Create(&itemDTOs).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderPersistent) mapToOrderDTO(order *entity.Order) *model.OrderDTO {
	return &model.OrderDTO{
		ID:                   order.ID,
		OrderNumber:          order.OrderNumber,
		UserID:               order.UserID,
		Status:               order.Status,
		Money:                order.Money,
		ShippingAddress:      order.ShippingAddress,
		ShippingPhone:        order.ShippingPhone,
		PaymentTransactionID: order.PaymentTransactionID,
		ShippingTrackingID:   order.ShippingTrackingID,
		CreatedAt:            order.CreateAt,
		UpdatedAt:            order.UpdateAt,
	}
}

func (r *orderPersistent) mapToOrderItemDTOs(items []entity.OrderItem, orderID string) []model.OrderItemDTO {
	if len(items) == 0 {
		return nil
	}

	dto := make([]model.OrderItemDTO, 0, len(items))
	for _, item := range items {
		dto = append(dto, model.OrderItemDTO{
			OrderID:     orderID,
			ID:          uuid.NewString(),
			UnitPrice:   item.UnitPrice,
			ProductID:   item.OrderID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
		})
	}
	return dto
}
