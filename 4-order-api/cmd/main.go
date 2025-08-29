package main

import (
	"demo/order-api/configs"
	"demo/order-api/internal/product"
	"demo/order-api/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	db := db.NewDB(config)

	// Register repositories
	productRepository := product.NewProductRepository(db.Db)

	// Register handlers
	product.NewProductHandler(router, product.ProductHandlerDeps{
		Config:            config,
		ProductRepository: productRepository,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: router,
	}

	fmt.Printf("Server is running on %s:%s", config.Server.Host, config.Server.Port)
	server.ListenAndServe()
}
