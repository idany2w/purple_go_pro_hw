package main

import (
	"demo/order-api/configs"
	"demo/order-api/internal/product"
	"demo/order-api/pkg/db"
)

func main() {
	config := configs.LoadConfig()
	db := db.NewDB(config)
	db.Migrate(&product.Product{})
}
