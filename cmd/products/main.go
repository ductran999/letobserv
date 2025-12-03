package main

import (
	"log"

	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/di"
)

func main() {
	config, err := configs.LoadProductConfig("product.env")
	if err != nil {
		log.Fatalln("failed to load product config:", err)
	}

	container, err := di.NewProductContainer(config)
	if err != nil {
		log.Fatalln("failed to create new container:", err)
	}

	if err := container.StartHTTPServer(); err != nil {
		log.Fatalln("start http server failed:", err)
	}
}
