package models

import (
	"database/sql"

	"github.com/dgrijalva/jwt-go"
)

type (
	User struct {
		Id        int    `json:"id"`
		UserName  string `json:"username"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	UserDb struct {
		Id        int
		UserName  string
		Password  []byte
		FirstName sql.NullString
		LastName  sql.NullString
		Email     sql.NullString
	}

	LoginCredentials struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	CustomClaims struct {
		UserName string `json:"username"`
		jwt.StandardClaims
	}
)
