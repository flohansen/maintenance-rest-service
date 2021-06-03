package models

import "github.com/dgrijalva/jwt-go"

type (
	LoginCredentials struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	CustomClaims struct {
		UserName string `json:"username"`
		jwt.StandardClaims
	}
)
