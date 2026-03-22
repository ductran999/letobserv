package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	loader "github.com/ductran999/letobserv/pkg/config"
	"github.com/ductran999/letobserv/services/inventory/internal/app"
	"github.com/ductran999/letobserv/services/inventory/internal/config"
)

func main() {
	cfg, err := loader.LoadConfig[config.Config](".env")
	if err != nil {
		log.Fatalln("load config failed:", err)
	}

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalln("init app failed:", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	if err := app.Run(ctx); err != nil {
		log.Fatalln("run app failed:", err)
	}
	stop()
}
