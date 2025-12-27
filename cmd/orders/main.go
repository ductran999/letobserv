package main

import (
	"log"

	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/bootstrap"
	"github.com/ductran999/letobserv/internal/transport"
)

func main() {
	config, err := configs.LoadOrderConfig()
	if err != nil {
		log.Fatalln("failed to load order service config:", err)
	}

	container, err := bootstrap.NewOrderContainer(config)
	if err != nil {
		log.Fatalln("failed to create new container:", err)
	}

	if err := transport.StartOrderHTTPServer(container); err != nil {
		log.Fatalln("start http server failed:", err)
	}
}
