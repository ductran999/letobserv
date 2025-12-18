package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ductran999/letobserv/internal/domain/repository"
	"github.com/ductran999/letobserv/pkg/httpclient"
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

func (uc *placeOrderUsecase) PlacePOrder(ctx context.Context) error {
	time.Sleep(200 * time.Millisecond)

	resp, err := uc.httpClient.Patch(ctx, "http://localhost:11000/api/products/2/reduce-stock", nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		return errors.New("failed")
	}

	return nil
}
