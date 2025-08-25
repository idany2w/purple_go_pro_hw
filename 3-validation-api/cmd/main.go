package main

import (
	"demo/validation/configs"
	"demo/validation/internal/email"
	"demo/validation/internal/verify"
	"demo/validation/internal/verify/storage"
	"fmt"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: router,
	}

	emailService := email.NewEmailService(config)
	hashStorage := storage.NewFileStorage()
	verifyService := verify.NewVerifyService(config, emailService, hashStorage)

	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		Config:        config,
		VerifyService: verifyService,
	})

	fmt.Println("Server is running on port", config.Server.Port)
	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
