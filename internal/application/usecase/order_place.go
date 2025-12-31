package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ductran999/letobserv/internal/application/inputs"
	"github.com/ductran999/letobserv/internal/application/outputs"
	"github.com/ductran999/letobserv/internal/consts"
	"github.com/ductran999/letobserv/internal/domain/entity"
	"github.com/ductran999/letobserv/internal/domain/repository"
	"github.com/ductran999/letobserv/pkg/httpclient"
	"github.com/google/uuid"
)

type placeOrderUsecase struct {
	orderRepo  repository.OrderRepository
	httpClient httpclient.Client
}

func NewOrderUseCase(orderRepo repository.OrderRepository, httpClient httpclient.Client) OrderUseCase {
	return &placeOrderUsecase{
		orderRepo:  orderRepo,
		httpClient: httpClient,
	}
}

func (uc *placeOrderUsecase) PlacePOrder(ctx context.Context, input inputs.PlaceOrderRequest) (*outputs.PlacedOrderOutput, error) {
	order := uc.mapPlaceOrderRequestToEntity(input)

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	resp := outputs.PlacedOrderOutput{
		ID:              order.ID,
		UserID:          order.UserID,
		OrderNumber:     order.OrderNumber,
		Money:           order.Money,
		ShippingAddress: order.ShippingAddress,
		ShippingPhone:   order.ShippingPhone,
	}

	return &resp, nil
}

func (uc *placeOrderUsecase) mapPlaceOrderRequestToEntity(input inputs.PlaceOrderRequest) *entity.Order {
	orderID := uuid.NewString()
	orderNumber := strings.ReplaceAll(orderID, "-", "")[:8]
	now := time.Now()

	items := make([]entity.OrderItem, 0, len(input.Items))
	for _, item := range input.Items {
		items = append(items, entity.OrderItem(item))
	}

	order := &entity.Order{
		UserID:               input.UserID,
		Money:                input.Money,
		ShippingAddress:      input.ShippingAddress,
		ShippingPhone:        input.ShippingPhone,
		PaymentTransactionID: uuid.NewString(),
		ShippingTrackingID:   uuid.NewString(),
		ID:                   orderID,
		OrderNumber:          fmt.Sprintf("ORD-%d-%s", time.Now().Year(), orderNumber),
		Status:               consts.OrderPendingStatus,
		Items:                items,
		CreateAt:             now,
		UpdateAt:             now,
	}

	return order
}
