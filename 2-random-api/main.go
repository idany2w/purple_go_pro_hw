package main

import (
	"demo/rand/handlers"
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	handlers.NewRandomHandler(router)

	fmt.Println("Server is running on port 8081")

	server.ListenAndServe()
}
