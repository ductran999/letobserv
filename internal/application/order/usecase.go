package order

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ductran999/letobserv/internal/domain/order"
	"github.com/ductran999/letobserv/pkg/httpclient"
	"github.com/google/uuid"
)

type OrderUseCase interface {
	PlacePOrder(ctx context.Context, input PlaceOrderRequest) (*PlacedOrderOutput, error)
}

type placeOrderUsecase struct {
	orderRepo  order.OrderRepository
	httpClient httpclient.Client
}

func NewOrderUseCase(orderRepo order.OrderRepository, httpClient httpclient.Client) OrderUseCase {
	return &placeOrderUsecase{
		orderRepo:  orderRepo,
		httpClient: httpClient,
	}
}

func (uc *placeOrderUsecase) PlacePOrder(ctx context.Context, input PlaceOrderRequest) (*PlacedOrderOutput, error) {
	order := uc.mapPlaceOrderRequestToEntity(input)

	if err := uc.orderRepo.Create(ctx, order); err != nil {
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

func (uc *placeOrderUsecase) mapPlaceOrderRequestToEntity(input PlaceOrderRequest) *order.Order {
	orderID := uuid.NewString()
	orderNumber := strings.ReplaceAll(orderID, "-", "")[:8]
	now := time.Now()

	items := make([]order.OrderItem, 0, len(input.Items))
	for _, item := range input.Items {
		items = append(items, order.OrderItem(item))
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
