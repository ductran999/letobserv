package repository

import (
	"context"

	"github.com/ductran999/letobserv/internal/domain/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
}
