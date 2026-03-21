package main

import (
	"fmt"
	"log"

	loader "github.com/ductran999/letobserv/pkg/config"
	"github.com/ductran999/letobserv/services/inventory/internal/config"
)

func main() {
	env, err := loader.LoadConfig[config.Config](".env")
	if err != nil {
		log.Fatalln("failed to load product config:", err)
	}
	fmt.Println(env)

	// app, err := bootstrap.NewInventory(env)
	// if err != nil {
	// 	log.Fatalln("failed to create new order app:", err)
	// }

	// if err := server.RunHTTP(env, app); err != nil {
	// 	log.Fatalln("start http server failed:", err)
	// }
}
