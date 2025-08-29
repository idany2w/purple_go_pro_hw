package main

import (
	"demo/order-api/configs"
	"demo/order-api/internal/auth"
	"demo/order-api/internal/product"
	"demo/order-api/internal/sms"
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

	// Register repositories
	productRepository := product.NewProductRepository(db.Db)

	// Register services
	smsService := sms.NewSmsService(db.Db)

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
