package auth

import (
	"demo/order-api/configs"
	"demo/order-api/internal/sms"
	"demo/order-api/pkg/jwt"
	"demo/order-api/pkg/request"
	"demo/order-api/pkg/response"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	SmsService *sms.SmsService
	JWT        *jwt.JWT
}

type AuthHandler struct {
	*configs.Config
	SmsService *sms.SmsService
	JWT        *jwt.JWT
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	authHandler := &AuthHandler{
		Config:     deps.Config,
		SmsService: deps.SmsService,
		JWT:        deps.JWT,
	}

	router.HandleFunc("POST /auth/login", authHandler.login())
	router.HandleFunc("POST /auth/register", authHandler.register())
}

func (h *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := request.HandleBody[sms.VerifySmsRequest](w, r)

		if body == nil {
			response.SendJsonError(&w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := h.SmsService.VerifySms(body.Phone, body.Token, body.Code)
		if err != nil {
			response.SendJsonError(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := h.JWT.Create(body.Phone)
		if err != nil {
			response.SendJsonError(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "Login successful", sms.VerifySmsResponse{Token: token}, http.StatusOK)
	}
}

func (h *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := request.HandleBody[sms.SendSmsRequest](w, r)

		if body == nil {
			response.SendJsonError(&w, "Invalid request body", http.StatusBadRequest)
			return
		}

		token, err := h.SmsService.SendSms(body.Phone)
		if err != nil {
			response.SendJsonError(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.SendJsonSuccess(&w, "SMS sent", sms.SendSmsResponse{Key: token}, http.StatusOK)
	}
}
