package main

import (
	"demo/validation/configs"
	"demo/validation/internal/email"
	"demo/validation/internal/verify"
	"fmt"
	"net/http"
)

func main() {
	config := configs.LoadConfig()
	router := http.NewServeMux()

	server := &http.Server{
		Addr:    config.Server.Port,
		Handler: router,
	}

	emailService := email.NewEmailService(config)
	verifyService := verify.NewVerifyService(emailService)

	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		Config:        config,
		VerifyService: verifyService,
	})

	fmt.Println("Server is running on port", config.Server.Port)
	server.ListenAndServe()
}
