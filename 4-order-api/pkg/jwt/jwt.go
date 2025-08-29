package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Key string
}

func NewJWT(key string) *JWT {
	return &JWT{
		Key: key,
	}
}

func (j *JWT) Create(phone string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": phone,
	})

	tokenString, err := token.SignedString([]byte(j.Key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWT) Validate(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Key), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}

	return claims["phone"].(string), nil
}
