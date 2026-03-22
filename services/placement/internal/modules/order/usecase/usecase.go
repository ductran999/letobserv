package ordersvc

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

type placeOrderUsecase struct {
	orderRepo       order.OrderRepository
	inventoryClient order.InventoryRestful
}

func NewOrderUseCase(
	orderRepo order.OrderRepository,
	inventoryClient order.InventoryRestful,
) order.UseCase {
	return &placeOrderUsecase{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
	}
}

func (uc *placeOrderUsecase) PlacePOrder(
	ctx context.Context,
	input order.PlaceOrderRequest,
) (*order.PlacedOrderOutput, error) {
	ord := uc.mapPlaceOrderRequestToEntity(input)
	if err := uc.orderRepo.Create(ctx, ord); err != nil {
		return nil, err
	}

	reserveReq := uc.prepareReserveRequest(ord)
	if err := uc.inventoryClient.Reserve(ctx, reserveReq); err != nil {
		return nil, err
	}

	resp := order.PlacedOrderOutput{
		ID:              ord.ID,
		UserID:          ord.UserID,
		OrderNumber:     ord.OrderNumber,
		Money:           ord.Money,
		ShippingAddress: ord.ShippingAddress,
		ShippingPhone:   ord.ShippingPhone,
	}

	return &resp, nil
}

func (uc *placeOrderUsecase) GetOrder(ctx context.Context, id string) (*order.PlacedOrder, error) {
	ord, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errs.Internal(err)
	}

	orderItems := make([]order.OrderItem, 0, len(ord.Items))
	for _, item := range ord.Items {
		_ = uc.inventoryClient.GetProduct(ctx, item.ProductID)

		orderItems = append(orderItems, order.OrderItem{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
		})
	}

	return &order.PlacedOrder{
		PlacedOrderOutput: order.PlacedOrderOutput{
			ID:              ord.ID,
			OrderNumber:     ord.OrderNumber,
			UserID:          ord.UserID,
			Money:           ord.Money,
			ShippingAddress: ord.ShippingAddress,
			ShippingPhone:   ord.ShippingPhone,
		},
		Items: orderItems,
	}, nil
}

func (uc *placeOrderUsecase) prepareReserveRequest(o *order.Order) order.InventoryReserveRequest {
	reserveItems := make([]order.ReserveItem, 0)
	for _, item := range o.Items {
		reserveItems = append(reserveItems, order.ReserveItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return order.InventoryReserveRequest{
		OrderID: o.ID,
		Items:   reserveItems,
		TTL:     defaultReserveTTL,
	}
}

func (uc *placeOrderUsecase) mapPlaceOrderRequestToEntity(input order.PlaceOrderRequest) *order.Order {
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
