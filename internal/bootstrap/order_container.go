package bootstrap

import (
	"fmt"

	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/application/usecase"
	"github.com/ductran999/letobserv/internal/domain/repository"
	"github.com/ductran999/letobserv/internal/infrastructure/persistent"
	"github.com/ductran999/letobserv/pkg/httpclient"
)

type OrderContainer struct {
	Env        *configs.OrdersConfigEnv
	httpClient httpclient.Client

	OrderUC   usecase.OrderUseCase
	OrderRepo repository.OrderRepository
}

func NewOrderContainer(env *configs.OrdersConfigEnv) (*OrderContainer, error) {
	pg, err := ConnectOrderDB(env)
	if err != nil {
		return nil, fmt.Errorf("failed to connect order db: %w", err)
	}

	client := httpclient.New()
	orderRepo := persistent.NewOrderRepo(pg)
	orderUC := usecase.NewOrderUseCase(orderRepo, client)

	container := &OrderContainer{
		Env:        env,
		httpClient: client,
		OrderUC:    orderUC,
		OrderRepo:  orderRepo,
	}

	return container, nil
}
