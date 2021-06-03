package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

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
	var res models.JsonResponse
	var user models.User

	// Read user from request body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		res.Code = 400
		res.Content = "Could not create new user"
		res.Send(w)
		return
	}

	// Insert the user into the database
	_, err = uc.Db.Query(
		"INSERT INTO users (username, password, firstName, lastName, email) VALUES (?, ?, ?, ?, ?)",
		user.UserName, user.Password, user.FirstName, user.LastName, user.Email,
	)

	if err != nil {
		res.Code = 400
		res.Content = "Internal error"
		res.Send(w)
		return
	}

	// Everything went fine.
	res.Code = 200
	res.Content = "Success"
	res.Send(w)
}

// Tries to login an user using the given login credentials and send a token if succeeded.
func (uc UserController) LoginUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// Deletes an user and removes it from the database.
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
