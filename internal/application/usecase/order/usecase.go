package order

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ductran999/letobserv/internal/application/port/service"
	"github.com/ductran999/letobserv/internal/domain/order"
	"github.com/google/uuid"
)

const (
	defaultReserveTTL = 300
)

type OrderUseCase interface {
	PlacePOrder(ctx context.Context, input PlaceOrderRequest) (*PlacedOrderOutput, error)
}

type placeOrderUsecase struct {
	orderRepo order.OrderRepository

	inventoryService service.InventoryRestful
}

func NewOrderUseCase(orderRepo order.OrderRepository, inventoryService service.InventoryRestful) OrderUseCase {
	return &placeOrderUsecase{
		orderRepo:        orderRepo,
		inventoryService: inventoryService,
	}
}

func (uc *placeOrderUsecase) PlacePOrder(ctx context.Context, input PlaceOrderRequest) (*PlacedOrderOutput, error) {
	order := uc.mapPlaceOrderRequestToEntity(input)
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	reserveReq := uc.prepareReserveRequest(order)
	if err := uc.inventoryService.Reserve(ctx, reserveReq); err != nil {
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

func (uc *placeOrderUsecase) prepareReserveRequest(order *order.Order) service.InventoryReserveRequest {
	reserveItems := make([]service.ReserveItem, 0)
	for _, item := range order.Items {
		reserveItems = append(reserveItems, service.ReserveItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return service.InventoryReserveRequest{
		OrderID: order.ID,
		Items:   reserveItems,
		TTL:     defaultReserveTTL,
	}
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
