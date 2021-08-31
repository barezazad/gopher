package model

import (
	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
	jwt.StandardClaims
}
