package model

import "github.com/dgrijalva/jwt-go"

type JwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
