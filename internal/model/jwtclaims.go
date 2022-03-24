package model

import (
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
	jwt.StandardClaims
}
