package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtData struct {
	Phone string
}

type JWT struct {
	Key string
}

func NewJWT(key string) *JWT {
	return &JWT{
		Key: key,
	}
}

func (j *JWT) Create(data JwtData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": data.Phone,
	})

	tokenString, err := token.SignedString([]byte(j.Key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWT) Parse(tokenString string) (bool, *JwtData) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Key), nil
	})

	if err != nil {
		return false, nil
	}

	phone := token.Claims.(jwt.MapClaims)["phone"]

	return token.Valid, &JwtData{
		Phone: phone.(string),
	}
}
