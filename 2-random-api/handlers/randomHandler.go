package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
)

type RandomHandler struct {
}

func NewRandomHandler(router *http.ServeMux) *RandomHandler {
	handler := &RandomHandler{}
	router.HandleFunc("/", handler.random())
	return handler
}

func (h *RandomHandler) random() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%d", rand.Intn(6)+1)))
	}
}
