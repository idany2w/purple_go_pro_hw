package response

import (
	"encoding/json"
	"net/http"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func SendJsonResponse(w *http.ResponseWriter, data Response, statusCode int) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(statusCode)
	json.NewEncoder(*w).Encode(data)
}

func SendJsonError(w *http.ResponseWriter, message string, statusCode int) {
	SendJsonResponse(w, Response{
		Status:  StatusError,
		Message: message,
	}, statusCode)
}

func SendJsonSuccess(w *http.ResponseWriter, message string, data any, statusCode int) {
	SendJsonResponse(w, Response{
		Status:  StatusSuccess,
		Message: message,
		Data:    data,
	}, statusCode)
}
