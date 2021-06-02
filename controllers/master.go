package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/kluddizz/maintenance-master/models"
)

type (
	MasterController struct {
		Db *sql.DB
	}
)

func NewMasterController(db *sql.DB) *MasterController {
	return &MasterController{
		Db: db,
	}
}

func (mc MasterController) GetMaster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}
	master := models.Master{}

	err := mc.Db.QueryRow(
		"SELECT id, name, host, port FROM masters WHERE id = ?",
		p.ByName("id"),
	).Scan(
		&master.Id, &master.Name, &master.Host, &master.Port,
	)

	if err != nil {
		res.Code = 400
		res.Content = fmt.Sprintf("Could not find master with id `%s`", p.ByName("id"))
	} else {
		res.Code = 200
		res.Content = master
	}

	rj, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	fmt.Fprintf(w, "%s", rj)
}

func (mc MasterController) CreateMaster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var m models.Master
	err := decoder.Decode(&m)
	if err != nil {
		panic(err)
	}

	_, err = mc.Db.Query(
		"INSERT INTO masters (name, host, port) VALUES (?, ?, ?)",
		m.Name, m.Host, m.Port,
	)

	if err != nil {
		panic(err.Error())
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "Success")
}
