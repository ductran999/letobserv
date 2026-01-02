package main

import (
	"log"

	"github.com/ductran999/letobserv/cmd/inventory/server"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/bootstrap"
)

func main() {
	env, err := configs.LoadProductConfig()
	if err != nil {
		log.Fatalln("failed to load product config:", err)
	}

	app, err := bootstrap.NewInventory(env)
	if err != nil {
		log.Fatalln("failed to create new order app:", err)
	}

	if err := server.RunHTTP(env, app); err != nil {
		log.Fatalln("start http server failed:", err)
	}
}
