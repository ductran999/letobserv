package repo

import (
	"context"

	"github.com/ductran999/letobserv/services/orders/port"
)

type orderRepo struct{}

func NewOrderRepo() port.OrderRepo {
	return &orderRepo{}
}

func (r *orderRepo) Create(ctx context.Context, productID int) error {
	return nil
}
