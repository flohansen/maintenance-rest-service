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
	r.GET("/master/:id", mc.GetMaster)
	r.POST("/master", mc.CreateMaster)

	http.ListenAndServe("localhost:3000", r)
}
