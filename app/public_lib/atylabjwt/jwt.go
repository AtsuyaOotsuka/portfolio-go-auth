package atylabjwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtSvcInterface interface {
	CreateJwt(config *JwtConfig) (string, error)
}

type JwtSvcStruct struct {
}

func NewJwtSvc() *JwtSvcStruct {
	return &JwtSvcStruct{}
}

type JwtConfig struct {
	Key   []byte
	Uuid  string
	Email string
	Exp   time.Time
}

func (s *JwtSvcStruct) CreateJwt(config *JwtConfig) (string, error) {
	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("user%s", config.Uuid),
		"email": config.Email,
		"exp":   config.Exp.Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
