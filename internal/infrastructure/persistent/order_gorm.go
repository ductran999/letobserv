package persistent

import (
	"context"

	"github.com/ductran999/letobserv/internal/domain/repository"
)

type orderPersistent struct{}

func NewOrderRepo() repository.OrderRepository {
	return &orderPersistent{}
}

func (r *orderPersistent) Create(ctx context.Context, productID int) error {
	return nil
}
