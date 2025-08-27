package main

import (
	"demo/order-api/configs"
	"demo/order-api/internal/order"
	"demo/order-api/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	db := db.NewDB(config)

	// Register repositories
	orderRepository := order.NewOrderRepository(db.Db)

	// Register handlers
	order.NewOrderHandler(router, order.OrderHandlerDeps{
		Config:          config,
		OrderRepository: orderRepository,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: router,
	}

	fmt.Printf("Server is running on %s:%s", config.Server.Host, config.Server.Port)
	server.ListenAndServe()
}
