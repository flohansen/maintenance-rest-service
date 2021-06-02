package main

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kluddizz/maintenance-master/config"
	"github.com/kluddizz/maintenance-master/controllers"
)

func main() {
	dbConfig, _ := config.ReadDatabaseConfig("./database.json")
	db, err := sql.Open("mysql", dbConfig.DataSourceName())
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	r := httprouter.New()

	uc := controllers.NewUserController()
	mc := controllers.NewMasterController(db)

	r.GET("/user/:id", uc.GetUser)

	r.GET("/masters", mc.GetMasters)
	r.POST("/masters", mc.CreateMaster)

	r.GET("/masters/:id", mc.GetMaster)
	r.PUT("/masters/:id", mc.UpdateMaster)
	r.DELETE("/masters/:id", mc.DeleteMaster)

	http.ListenAndServe("localhost:3000", r)
}
