package usecase

import (
	"context"

	"github.com/ductran999/letobserv/internal/application/inputs"
	"github.com/ductran999/letobserv/internal/application/outputs"
)

type OrderUseCase interface {
	PlacePOrder(ctx context.Context, input inputs.PlaceOrderRequest) (*outputs.PlacedOrderOutput, error)
}
