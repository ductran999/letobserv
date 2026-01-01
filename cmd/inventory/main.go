package main

import (
	"log"

	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/bootstrap"
)

func main() {
	config, err := configs.LoadProductConfig()
	if err != nil {
		log.Fatalln("failed to load product config:", err)
	}

	container, err := bootstrap.NewInventory(config)
	if err != nil {
		log.Fatalln("failed to create new container:", err)
	}

	if err := container.StartHTTPServer(); err != nil {
		log.Fatalln("start http server failed:", err)
	}
}
