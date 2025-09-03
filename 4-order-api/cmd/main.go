package main

import (
	"demo/order-api/configs"
	"demo/order-api/internal/auth"
	"demo/order-api/internal/order"
	"demo/order-api/internal/product"
	"demo/order-api/internal/sms"
	"demo/order-api/migrations"
	"demo/order-api/pkg/db"
	"demo/order-api/pkg/jwt"
	"demo/order-api/pkg/middleware"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	db := db.NewDB(config)

	// Выполняем автоматические миграции
	err := migrations.AutoMigrate(db.Db)
	if err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	// Создаем тестовые данные
	err = migrations.CreateSampleData(db.Db)
	if err != nil {
		logrus.Warnf("Failed to create sample data: %v", err)
	}

	// Register repositories
	productRepository := product.NewProductRepository(db.Db)
	orderRepository := order.NewOrderRepository(db.Db)

	// Register services
	smsService := sms.NewSmsService(db.Db)
	orderService := order.NewOrderService(orderRepository)

	// Register handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:     config,
		SmsService: smsService,
		JWT:        jwt.NewJWT(config.Jwt.Key),
	})

	product.NewProductHandler(router, product.ProductHandlerDeps{
		Config:            config,
		ProductRepository: productRepository,
	})

	order.NewOrderHandler(router, order.OrderHandlerDeps{
		OrderService: orderService,
		JWT:          jwt.NewJWT(config.Jwt.Key),
	})

	// Register middleware

	logger := initLogger()
	deps := middleware.LoggingDeps{
		Logger: logger,
	}

	stack := middleware.Chain(
		middleware.Logging(&deps),
	)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: stack(router),
	}

	fmt.Printf("Server is running on %s:%s\n", config.Server.Host, config.Server.Port)
	server.ListenAndServe()
}

func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	return logger
}
