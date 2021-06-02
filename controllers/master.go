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

func (mc MasterController) GetMasters(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}

	var masters []models.Master

	query, err := mc.Db.Query("SELECT id, name, host, port FROM masters")
	defer query.Close()

	if err != nil {
		res.Code = 400
		res.Content = "Could not find masters"
		res.Send(w)
		return
	}

	for query.Next() {
		var master models.Master
		err := query.Scan(&master.Id, &master.Name, &master.Host, &master.Port)

		if err != nil {
			res.Code = 400
			res.Content = "Internal error"
			res.Send(w)
			return
		}

		masters = append(masters, master)
	}

	res.Code = 200
	res.Content = masters
	res.Send(w)
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
		res.Send(w)
		return
	}

	res.Code = 200
	res.Content = master
	res.Send(w)
}

func (mc MasterController) CreateMaster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}
	decoder := json.NewDecoder(r.Body)

	var m models.Master
	err := decoder.Decode(&m)

	if err != nil {
		res.Code = 400
		res.Content = "Could not parse json body"
		res.Send(w)
		return
	}

	_, err = mc.Db.Query(
		"INSERT INTO masters (name, host, port) VALUES (?, ?, ?)",
		m.Name, m.Host, m.Port,
	)

	if err != nil {
		res.Code = 400
		res.Content = "Could not create new master. Please check if the name is unique."
		res.Send(w)
		return
	}

	res.Code = 200
	res.Content = "Success"
	res.Send(w)
}

func (mc MasterController) UpdateMaster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}
	decoder := json.NewDecoder(r.Body)

	var m models.Master
	err := decoder.Decode(&m)

	if err != nil {
		res.Code = 400
		res.Content = "Could not parse json body"
		res.Send(w)
		return
	}

	_, err = mc.Db.Query(
		"UPDATE masters SET name = ?, host = ?, port = ? WHERE id = ?",
		m.Name, m.Host, m.Port, p.ByName("id"),
	)

	if err != nil {
		res.Code = 400
		res.Content = "Could not update master"
		res.Send(w)
		return
	}

	res.Code = 200
	res.Content = "Success"
	res.Send(w)
}

func (mc MasterController) DeleteMaster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}

	_, err := mc.Db.Query(
		"DELETE FROM masters WHERE id = ?",
		p.ByName("id"),
	)

	if err != nil {
		res.Code = 400
		res.Content = fmt.Sprintf("Could not delete master with id `%s`", p.ByName("id"))
		res.Send(w)
		return
	}

	res.Code = 200
	res.Content = "Success"
	res.Send(w)
}
