package main

import (
	"log"

	"github.com/ductran999/letobserv/cmd/orders/server"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/internal/bootstrap"
)

func main() {
	env, err := configs.LoadOrderConfig()
	if err != nil {
		log.Fatalln("failed to load order service config:", err)
	}

	app, err := bootstrap.NewOrderBootstrap(env)
	if err != nil {
		log.Fatalln("failed to create new container:", err)
	}

	if err := server.RunHTTP(app, env); err != nil {
		if env.ApmEnable {
			app.APMAgent.Shutdown()
		}
		log.Fatalln("http server error:", err)
	}
}
