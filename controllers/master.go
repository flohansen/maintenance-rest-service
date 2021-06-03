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
	MasterController struct {
		Db *sql.DB
	}
)

// Creates a new master controller, which handles interaction with the database.
func NewMasterController(db *sql.DB) *MasterController {
	return &MasterController{
		Db: db,
	}
}

// Requests all masters stored in the database.
func (mc MasterController) GetMasters(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}
	var masters []models.Master

	// Send the select query to the database to fetch stored master endpoints.
	query, err := mc.Db.Query("SELECT id, name, host, port FROM masters")
	defer query.Close()

	if err != nil {
		res.Code = 400
		res.Content = "Could not find masters"
		res.Send(w)
		return
	}

	// Fill the list of masters by iterating over all requested rows.
	for query.Next() {
		// Create and fill a new master object.
		var master models.Master
		err := query.Scan(&master.Id, &master.Name, &master.Host, &master.Port)

		if err != nil {
			res.Code = 400
			res.Content = "Internal error"
			res.Send(w)
			return
		}

		// Append the object to the array.
		masters = append(masters, master)
	}

	// Everything was successfull.
	res.Code = 200
	res.Content = masters
	res.Send(w)
}

// Requests a specific master identified by an id.
func (mc MasterController) GetMaster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}
	master := models.Master{}

	// Request the master endpoint using the given id from the database.
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

	// Everything went fine.
	res.Code = 200
	res.Content = master
	res.Send(w)
}

// Creates a new master object and stores it into the database.
func (mc MasterController) CreateMaster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}
	decoder := json.NewDecoder(r.Body)

	// Read the master object from the JSON body.
	var m models.Master
	err := decoder.Decode(&m)

	if err != nil {
		res.Code = 400
		res.Content = "Could not parse json body"
		res.Send(w)
		return
	}

	// Store the master object into the database.
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

	// Everything went fine.
	res.Code = 200
	res.Content = "Success"
	res.Send(w)
}

// Updates a master object stored in the database.
func (mc MasterController) UpdateMaster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}
	decoder := json.NewDecoder(r.Body)

	// Decode the master object from the JSON body.
	var m models.Master
	err := decoder.Decode(&m)

	if err != nil {
		res.Code = 400
		res.Content = "Could not parse json body"
		res.Send(w)
		return
	}

	// Update the master object inside the database using the decoded object.
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

	// Everything went fine.
	res.Code = 200
	res.Content = "Success"
	res.Send(w)
}

// Deletes a master from the database.
func (mc MasterController) DeleteMaster(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res := models.JsonResponse{}

	// Remove the master object from the database using the id.
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

	// Everything went fine.
	res.Code = 200
	res.Content = "Success"
	res.Send(w)
}
