package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/ductran999/letobserv/pkg/apm"
	"github.com/ductran999/letobserv/pkg/database"
	"github.com/ductran999/letobserv/pkg/httpclient"
	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/ductran999/letobserv/pkg/middleware"
	gen "github.com/ductran999/letobserv/services/placement/api/gen/orders"
	configs "github.com/ductran999/letobserv/services/placement/internal/config"
	orderrepo "github.com/ductran999/letobserv/services/placement/internal/modules/order/infrastructure/repository"
	"github.com/ductran999/letobserv/services/placement/internal/modules/order/infrastructure/service"
	orderhttp "github.com/ductran999/letobserv/services/placement/internal/modules/order/transport/order"
	ordersvc "github.com/ductran999/letobserv/services/placement/internal/modules/order/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type PlacementApp struct {
	cfg        *configs.Config
	APMAgent   apm.APMAgent
	logger     *zerolog.Logger
	httpClient httpclient.Client

	apiServer *http.Server
}

type server struct {
	orderhttp.OrderHandler
}

func NewPlacementApp(cfg *configs.Config) (*PlacementApp, error) {
	pg := database.MustConnect(cfg.PgConfig())
	slog.Info("connect pg db successfully!")

	apmAgent, err := apm.NewAgentIfEnabled(context.Background(), cfg.APMConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to start apm agent: %w", err)
	}
	if apmAgent != nil {
		slog.Info("apm agent started successfully!")
	}

	logger, err := logger.New(cfg.GetLoggerConfig(), logger.WithAPM(cfg.GetLoggerAPMConfig()))
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	slog.Info("initialize logger successfully!")

	if cfg.ServiceEnv == configs.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	client := httpclient.New()
	orderRepo := orderrepo.NewOrderRepo(pg)
	inventorySvc := service.NewInventoryService(client, cfg.InventoryServiceBaseURL)
	orderUC := ordersvc.NewOrderUseCase(orderRepo, inventorySvc)
	orderHandler := orderhttp.NewOrderHandler(orderUC)

	router := gin.Default()
	if apmAgent != nil {
		router.Use(otelgin.Middleware(cfg.ServiceName))
	}
	router.Use(middleware.LoggingMiddleware())

	gen.RegisterHandlers(router, &server{
		OrderHandler: orderHandler,
	})

	app := &PlacementApp{
		cfg:        cfg,
		httpClient: client,
		logger:     logger,
		apiServer: &http.Server{
			Addr:              net.JoinHostPort(cfg.ServiceHost, cfg.ServicePort),
			Handler:           router,
			ReadHeaderTimeout: 500 * time.Millisecond,
		},
	}

	return app, nil
}

func (a *PlacementApp) Run(ctx context.Context) error {
	go func() {
		slog.Info("api server started", "addr", a.apiServer.Addr)
		if err := a.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("api server error", "err", err)
		}
	}()

	<-ctx.Done()

	return a.shutdown(ctx)
}

func (a *PlacementApp) shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
	defer cancel()

	if a.apiServer != nil {
		if err := a.apiServer.Shutdown(shutdownCtx); err != nil {
			slog.Warn("shutdown api server failed", "err", err)
		}
	}

	if a.APMAgent != nil {
		if err := a.APMAgent.Shutdown(shutdownCtx); err != nil {
			slog.Warn("shutdown apm agent failed", "err", err)
		}
	}

	return nil
}
