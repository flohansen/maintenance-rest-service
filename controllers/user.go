package controllers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/kluddizz/maintenance-rest-service/models"
)

type (
	UserController struct {
		Db *sql.DB
	}
)

// Creates a new instance of the user controller structure.
func NewUserController(db *sql.DB) *UserController {
	return &UserController{
		Db: db,
	}
}

// Creates a new user and add it to the database.
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.NewJsonResponse(w)
	var user models.User

	// Read user from request body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		res.Code = 400
		res.Content = "Internal error"
		res.Send()
		return
	}

	// Insert the user into the database
	_, err = uc.Db.Query(
		"INSERT INTO users (username, password, firstName, lastName, email) VALUES (?, ?, ?, ?, ?)",
		user.UserName, user.Password, user.FirstName, user.LastName, user.Email,
	)

	if err != nil {
		res.Code = 400
		res.Content = "Could not create new user"
		res.Send()
		return
	}

	// Everything went fine.
	res.Code = 200
	res.Content = "Success"
	res.Send()
}

// Tries to login an user using the given login credentials and send a token if succeeded.
func (uc UserController) LoginUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.NewJsonResponse(w)
	var loginCredentials models.LoginCredentials
	privateKey, err := ioutil.ReadFile("./private.key")

	if err != nil {
		res.Code = 400
		res.Content = "Internal error"
		res.Send()
		return
	}

	// Read credentials from request body
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&loginCredentials)

	if err != nil {
		res.Code = 400
		res.Content = "Internal error"
		res.Send()
		return
	}

	// Check if the user exists
	err = uc.Db.QueryRow(
		"SELECT * FROM users WHERE username = ? AND password = ?",
		loginCredentials.UserName, loginCredentials.Password,
	).Err()

	if err != nil {
		res.Code = 400
		res.Content = "Could not find user"
		res.Send()
		return
	}

	// Create token
	claims := models.CustomClaims{
		UserName: loginCredentials.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer:    "maintenance-rest-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(privateKey)

	if err != nil {
		res.Code = 400
		res.Content = "Internal error"
		res.Send()
		return
	}

	// Everything went fine.
	res.Code = 200
	res.Content = signedToken
	res.Response.Header().Set("Authorization", "Bearer "+signedToken)
	res.Send()
}

// Deletes an user and removes it from the database.
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.NewJsonResponse(w)
	var loginCredentials models.LoginCredentials

	// Parse user from the request body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginCredentials)

	if err != nil {
		res.Code = 400
		res.Content = "Internal error"
		res.Send()
		return
	}

	// Remove user from database
	_, err = uc.Db.Query(
		"DELETE FROM users WHERE username = ? AND password = ?",
		loginCredentials.UserName, loginCredentials.Password,
	)

	if err != nil {
		res.Code = 400
		res.Content = "Could not delete user"
		res.Send()
		return
	}

	// Everything went fine
	res.Code = 200
	res.Content = "Success"
	res.Send()
}
