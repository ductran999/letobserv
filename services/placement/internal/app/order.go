package app

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/ductran999/dbkit"
// 	"github.com/ductran999/letobserv/configs"
// 	"github.com/ductran999/letobserv/pkg/apm"
// 	"github.com/ductran999/letobserv/pkg/httpclient"
// 	"github.com/ductran999/letobserv/pkg/logger"
// 	"gorm.io/gorm"
// )

// type OrderBootstrap struct {
// 	OrderUC  order.OrderUseCase
// 	APMAgent apm.APMAgent

// 	env        *configs.OrdersConfigEnv
// 	httpClient httpclient.Client
// 	orderRepo  domainOrder.OrderRepository
// }

// func NewOrderBootstrap(env *configs.OrdersConfigEnv) (*OrderBootstrap, error) {
// 	pg, err := connectOrderDB(env)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect order db: %w", err)
// 	}
// 	log.Println("[INFO] connect pg db successfully!")

// 	var apmAgent apm.APMAgent
// 	if env.ApmEnable {
// 		apmAgent, err = newOrderServiceAPMAgent(env)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to start apm agent: %w", err)
// 		}
// 		log.Println("[INFO] start apm agent successfully!")
// 	}

// 	if err = newOrdersLogger(env); err != nil {
// 		return nil, fmt.Errorf("init logger: %w", err)
// 	}
// 	log.Println("[INFO] initialize logger successfully!")

// 	client := httpclient.New()
// 	orderRepo := infraOrder.NewOrderRepo(pg)
// 	inventorySvc := infraService.NewInventoryService(client, env.InventoryServiceBaseURL)
// 	orderUC := order.NewOrderUseCase(orderRepo, inventorySvc)

// 	app := &OrderBootstrap{
// 		env:        env,
// 		httpClient: client,
// 		OrderUC:    orderUC,
// 		orderRepo:  orderRepo,
// 		APMAgent:   apmAgent,
// 	}

// 	return app, nil
// }

// func connectOrderDB(env *configs.OrdersConfigEnv) (*gorm.DB, error) {
// 	pgConf := dbkit.PostgreSQLConfig{
// 		Config: dbkit.Config{
// 			Host:     env.PgHost,
// 			Port:     env.PgPort,
// 			Username: env.PgUser,
// 			Password: env.PgPass,
// 			Database: env.PgDB,
// 			TimeZone: "UTC",
// 		},
// 		SSLMode: "disable",
// 	}

// 	return dbconn.ConnectDB(pgConf)
// }

// func newOrderServiceAPMAgent(env *configs.OrdersConfigEnv) (apm.APMAgent, error) {
// 	config := apm.AgentConfig{
// 		ServiceInfo: apm.ServiceInfo{
// 			Name:    env.ServiceName,
// 			Version: env.ServiceVersion,
// 			Env:     env.ServiceEnv,
// 		},
// 		APIKey:           env.ApmApiKey,
// 		InsecureMode:     env.ApmInsecureMode,
// 		ExporterEndpoint: env.ApmExporterEndpoint,
// 		Enabled:          env.ApmEnable,
// 	}
// 	agent, err := apm.NewAgent(context.Background(), config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return agent, nil
// }

// func newOrdersLogger(env *configs.OrdersConfigEnv) error {
// 	serviceInfo := logger.ServiceInfo{
// 		Name:    env.ServiceName,
// 		Version: env.ServiceVersion,
// 		Env:     env.ServiceEnv,
// 	}

// 	if !env.ApmEnable {
// 		return logger.New(serviceInfo)
// 	}

// 	apmConfig := logger.APMConfig{
// 		Enable:           env.ApmEnable,
// 		ExporterEndpoint: env.ApmExporterEndpoint,
// 		ApiKey:           env.ApmApiKey,
// 	}

// 	return logger.New(serviceInfo, logger.WithAPM(apmConfig))
// }
