package bootstrap

import (
	"fmt"

	"github.com/ductran999/dbkit"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/application/order"
	domainOrder "github.com/ductran999/letobserv/internal/domain/order"
	infraOrder "github.com/ductran999/letobserv/internal/infrastructure/order"
	"github.com/ductran999/letobserv/pkg/dbconn"
	"github.com/ductran999/letobserv/pkg/httpclient"
	"gorm.io/gorm"
)

type OrderBootstrap struct {
	OrderUC order.OrderUseCase

	env        *configs.OrdersConfigEnv
	httpClient httpclient.Client
	orderRepo  domainOrder.OrderRepository
}

func NewOrderBootstrap(env *configs.OrdersConfigEnv) (*OrderBootstrap, error) {
	pg, err := connectOrderDB(env)
	if err != nil {
		return nil, fmt.Errorf("failed to connect order db: %w", err)
	}

	client := httpclient.New()
	orderRepo := infraOrder.NewOrderRepo(pg)
	orderUC := order.NewOrderUseCase(orderRepo, client)

	app := &OrderBootstrap{
		env:        env,
		httpClient: client,
		OrderUC:    orderUC,
		orderRepo:  orderRepo,
	}

	return app, nil
}

func connectOrderDB(env *configs.OrdersConfigEnv) (*gorm.DB, error) {
	pgConf := dbkit.PostgreSQLConfig{
		Config: dbkit.Config{
			Host:     env.PgHost,
			Port:     env.PgPort,
			Username: env.PgUser,
			Password: env.PgPass,
			Database: env.PgDB,
			TimeZone: "UTC",
		},
		SSLMode: "disable",
	}

	return dbconn.ConnectDB(pgConf)
}
