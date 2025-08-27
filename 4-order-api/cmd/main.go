package main

import (
	"demo/order-api/configs"
	"demo/order-api/pkg/db"
	"fmt"
)

func main() {
	config := configs.LoadConfig()
	db := db.NewDB(config)
	fmt.Println(db)
}
