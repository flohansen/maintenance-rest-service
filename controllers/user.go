package controllers

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

}

// Tries to login an user using the given login credentials and send a token if succeeded.
func (uc UserController) LoginUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// Deletes an user and removes it from the database.
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
