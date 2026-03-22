package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ductran999/letobserv/pkg/apm"
	"github.com/ductran999/letobserv/pkg/database"
	"github.com/ductran999/letobserv/pkg/logger"
	"github.com/ductran999/letobserv/pkg/txmanager"
	gen "github.com/ductran999/letobserv/services/inventory/api/gen/openapi"
	"github.com/ductran999/letobserv/services/inventory/internal/config"
	"github.com/ductran999/letobserv/services/inventory/internal/middleware"
	commonhttp "github.com/ductran999/letobserv/services/inventory/internal/modules/common/delivery/http"
	inventoryhttp "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/delivery/http/handler"
	inventoryrepo "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/infra/repository"
	inventoryuc "github.com/ductran999/letobserv/services/inventory/internal/modules/inventory/usecase"
	producthttp "github.com/ductran999/letobserv/services/inventory/internal/modules/product/delivery/http/handler"
	productrepo "github.com/ductran999/letobserv/services/inventory/internal/modules/product/infra/repository"
	productsvc "github.com/ductran999/letobserv/services/inventory/internal/modules/product/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/gorm"
)

type InventoryApp struct {
	cfg            *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	APMAgent       apm.APMAgent
	apiServer      *http.Server
	internalServer *http.Server
}

type server struct {
	producthttp.ProductHandler
	inventoryhttp.InventoryHandler
}

func NewApp(cfg *config.Config) (*InventoryApp, error) {
	db := database.MustConnect(cfg.PgConfig())
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

	if cfg.ServiceEnv == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	if apmAgent != nil {
		router.Use(otelgin.Middleware(cfg.ServiceName))
	}
	router.Use(middleware.LoggingMiddleware())

	productRepo := productrepo.NewRepository(db)
	productSvc := productsvc.New(productRepo)
	productHdl := producthttp.New(productSvc)
	inventoryRepo := inventoryrepo.NewInventoryRepo(db)
	inventorySvc := inventoryuc.NewInventoryUsecase(txmanager.NewTransactionManager(db), productRepo, inventoryRepo)
	inventoryHdl := inventoryhttp.NewHandler(inventorySvc)

	gen.RegisterHandlers(router, &server{
		ProductHandler:   productHdl,
		InventoryHandler: inventoryHdl,
	})

	internalRouter := gin.New()
	commonHdl := commonhttp.New(db)
	internalRouter.GET("/health", commonHdl.HealthCheck)
	internalRouter.GET("/ready", commonHdl.ReadinessCheck)

	return &InventoryApp{
		cfg:      cfg,
		db:       db,
		logger:   logger,
		APMAgent: apmAgent,
		apiServer: &http.Server{
			Addr:    net.JoinHostPort(cfg.ServiceHost, cfg.ServicePort),
			Handler: router,
		},
		internalServer: &http.Server{
			Addr:    net.JoinHostPort("0.0.0.0", "8081"),
			Handler: internalRouter,
		},
	}, nil
}

func (a *InventoryApp) Run() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("api server started", "addr", a.apiServer.Addr)
		if err := a.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("api server error", "err", err)
		}
	}()

	go func() {
		slog.Info("internal server started", "addr", a.internalServer.Addr)
		if err := a.internalServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("internal server error", "err", err)
		}
	}()

	<-quit
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.internalServer.Shutdown(ctx)
	return a.apiServer.Shutdown(ctx)
}
