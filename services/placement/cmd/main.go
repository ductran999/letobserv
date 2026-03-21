package main

import (
	"log"

	loader "github.com/ductran999/letobserv/pkg/config"
	configs "github.com/ductran999/letobserv/services/placement/internal/config"
)

func main() {
	env, err := loader.LoadConfig[configs.Config](".env")
	if err != nil {
		log.Fatalln("failed to load order service config:", err)
	}
	log.Println(env)

	// app, err := bootstrap.NewOrderBootstrap(env)
	// if err != nil {
	// 	log.Fatalln("failed to create new container:", err)
	// }

	// if err := server.RunHTTP(app, env); err != nil {
	// 	if env.ApmEnable {
	// 		if err := app.APMAgent.Shutdown(context.Background()); err != nil {
	// 			log.Println("shutdown apm error: %w", err)
	// 		}
	// 	}
	// 	log.Fatalln("http server error:", err)
	// }
}
