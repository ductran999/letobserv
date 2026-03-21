package orderuc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ductran999/letobserv/pkg/errs"
	order "github.com/ductran999/letobserv/services/placement/internal/modules/order/domain"
	"github.com/google/uuid"
)

const (
	defaultReserveTTL = 300
)

type OrderUseCase interface {
	PlacePOrder(ctx context.Context, input PlaceOrderRequest) (*PlacedOrderOutput, error)

	GetOrder(ctx context.Context, id string) (*PlacedOrder, error)
}

type placeOrderUsecase struct {
	orderRepo      order.OrderRepository
	inventoryorder order.InventoryRestful
}

func NewOrderUseCase(orderRepo order.OrderRepository, inventoryorder order.InventoryRestful) OrderUseCase {
	return &placeOrderUsecase{
		orderRepo:      orderRepo,
		inventoryorder: inventoryorder,
	}
}

func (uc *placeOrderUsecase) PlacePOrder(ctx context.Context, input PlaceOrderRequest) (*PlacedOrderOutput, error) {
	order := uc.mapPlaceOrderRequestToEntity(input)
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	reserveReq := uc.prepareReserveRequest(order)
	if err := uc.inventoryorder.Reserve(ctx, reserveReq); err != nil {
		return nil, err
	}

	resp := PlacedOrderOutput{
		ID:              order.ID,
		UserID:          order.UserID,
		OrderNumber:     order.OrderNumber,
		Money:           order.Money,
		ShippingAddress: order.ShippingAddress,
		ShippingPhone:   order.ShippingPhone,
	}

	return &resp, nil
}

func (uc *placeOrderUsecase) GetOrder(ctx context.Context, id string) (*PlacedOrder, error) {
	order, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errs.Internal(err)
	}

	orderItems := make([]OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		_ = uc.inventoryorder.GetProduct(ctx, item.ProductID)

		orderItems = append(orderItems, OrderItem{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
		})
	}

	return &PlacedOrder{
		PlacedOrderOutput: PlacedOrderOutput{
			ID:              order.ID,
			OrderNumber:     order.OrderNumber,
			UserID:          order.UserID,
			Money:           order.Money,
			ShippingAddress: order.ShippingAddress,
			ShippingPhone:   order.ShippingPhone,
		},
		Items: orderItems,
	}, nil
}

func (uc *placeOrderUsecase) prepareReserveRequest(o *order.Order) order.InventoryReserveRequest {
	// reserveItems := make([]order.ReserveItem, 0)
	// for _, item := range order.Items {
	// 	reserveItems = append(reserveItems, order.ReserveItem{
	// 		ProductID: item.ProductID,
	// 		Quantity:  item.Quantity,
	// 	})
	// }

	// return order.InventoryReserveRequest{
	// 	OrderID: order.ID,
	// 	Items:   reserveItems,
	// 	TTL:     defaultReserveTTL,
	// }
	return order.InventoryReserveRequest{}
}

func (uc *placeOrderUsecase) mapPlaceOrderRequestToEntity(input PlaceOrderRequest) *order.Order {
	orderID := uuid.NewString()
	orderNumber := strings.ReplaceAll(orderID, "-", "")[:8]
	now := time.Now()

	items := make([]order.OrderItem, 0, len(input.Items))
	for _, item := range input.Items {
		orderItem := order.OrderItem{
			ID:          uuid.NewString(),
			OrderID:     orderID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
		items = append(items, orderItem)
	}

	order := &order.Order{
		UserID:               input.UserID,
		Money:                input.Money,
		ShippingAddress:      input.ShippingAddress,
		ShippingPhone:        input.ShippingPhone,
		PaymentTransactionID: uuid.NewString(),
		ShippingTrackingID:   uuid.NewString(),
		ID:                   orderID,
		OrderNumber:          fmt.Sprintf("ORD-%d-%s", time.Now().Year(), orderNumber),
		Status:               order.OrderPendingStatus,
		Items:                items,
		CreateAt:             now,
		UpdateAt:             now,
	}

	return order
}
