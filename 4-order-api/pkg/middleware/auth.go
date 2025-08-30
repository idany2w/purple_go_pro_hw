package middleware

import (
	"context"
	"net/http"
	"strings"

	"demo/order-api/configs"
	"demo/order-api/pkg/jwt"
	"demo/order-api/pkg/response"
)

type key string

const (
	ContextPhoneKey key = "ContextPhoneKey"
)

func Auth(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			Fail(&w)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		valid, data := jwt.NewJWT(config.Jwt.Key).Parse(token)

		if !valid {
			Fail(&w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextPhoneKey, data.Phone)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func Fail(w *http.ResponseWriter) {
	response.SendJsonError(w, "Unauthorized", http.StatusUnauthorized)
}
