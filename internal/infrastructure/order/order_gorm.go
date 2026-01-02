package order

import (
	"context"

	"github.com/ductran999/letobserv/internal/domain/order"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderPersistent struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) order.OrderRepository {
	return &orderPersistent{db: db}
}

func (r *orderPersistent) Create(ctx context.Context, order *order.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		orderDTO := r.mapToOrderDTO(order)
		if err := tx.WithContext(ctx).Table(orderDTO.TableName()).Create(orderDTO).Error; err != nil {
			return err
		}

		itemDTOs := r.mapToOrderItemDTOs(order.Items, orderDTO.ID)
		if err := tx.WithContext(ctx).Table((&OrderItemDTO{}).TableName()).Create(&itemDTOs).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderPersistent) mapToOrderDTO(order *order.Order) *OrderDTO {
	return &OrderDTO{
		ID:                   order.ID,
		OrderNumber:          order.OrderNumber,
		UserID:               order.UserID,
		Status:               string(order.Status),
		Money:                order.Money,
		ShippingAddress:      order.ShippingAddress,
		ShippingPhone:        order.ShippingPhone,
		PaymentTransactionID: order.PaymentTransactionID,
		ShippingTrackingID:   order.ShippingTrackingID,
		CreatedAt:            order.CreateAt,
		UpdatedAt:            order.UpdateAt,
	}
}

func (r *orderPersistent) mapToOrderItemDTOs(items []order.OrderItem, orderID string) []OrderItemDTO {
	if len(items) == 0 {
		return nil
	}

	dto := make([]OrderItemDTO, 0, len(items))
	for _, item := range items {
		dto = append(dto, OrderItemDTO{
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
