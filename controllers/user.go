package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kluddizz/maintenance-rest-service/models"
)

type (
	UserController struct {
		Db *sql.DB
	}
)

func NewUserController(db *sql.DB) *UserController {
	return &UserController{
		Db: db,
	}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	u := models.User{
		Id:        p.ByName("id"),
		FirstName: "Florian",
		LastName:  "Hansen",
	}

	uj, _ := json.Marshal(u)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}
