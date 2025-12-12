package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ductran999/letobserv/pkg/httpclient"
	"github.com/ductran999/letobserv/services/orders/port"
)

type OrderUseCase struct {
	repo       port.OrderRepo
	httpClient httpclient.Client
}

func NewOrderUseCase(repo port.OrderRepo) port.OrderUseCase {
	return &OrderUseCase{
		repo:       repo,
		httpClient: httpclient.New(),
	}
}

func (uc *OrderUseCase) PlacePOrder(ctx context.Context) error {
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
