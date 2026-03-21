package app

// func RunHTTP(env *configs.InventoryEnv, app *bootstrap.InventoryBootstrap) error {
// 	if env.ServiceEnv == consts.ProductionEnv {
// 		gin.SetMode(gin.ReleaseMode)
// 	}
// 	hdl := inventory.NewInventoryHandler(app.InventoryUC)

// 	r := gin.Default()
// 	if env.ApmEnable {
// 		r.Use(otelgin.Middleware(env.ServiceName))
// 	}
// 	r.Use(middleware.LoggingMiddleware())
// 	generated.RegisterHandlers(r, hdl)

// 	// Start http server
// 	address := net.JoinHostPort("0.0.0.0", env.ServicePort)
// 	log.Println("[INFO] product service serving on", address)
// 	if err := r.Run(address); err != nil && errors.Is(err, http.ErrServerClosed) {
// 		return err
// 	}

// 	return nil
// }
