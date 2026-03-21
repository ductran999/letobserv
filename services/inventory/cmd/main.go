package main

import (
	"log"
	"log/slog"

	loader "github.com/ductran999/letobserv/pkg/config"
	"github.com/ductran999/letobserv/services/inventory/internal/app"
	"github.com/ductran999/letobserv/services/inventory/internal/config"
)

func main() {
	cfg, err := loader.LoadConfig[config.Config](".env")
	if err != nil {
		log.Fatalln("load config failed", err)
	}
	slog.Info("load config successfully!")

	_, err = app.NewApp(cfg)
	if err != nil {
		log.Fatalln("init app failed:", err)
	}
	slog.Info("app initialized successfully!")
}
