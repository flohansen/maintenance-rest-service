package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/kluddizz/maintenance-rest-service/config"
	"github.com/kluddizz/maintenance-rest-service/controllers"
	"github.com/kluddizz/maintenance-rest-service/middlewares"
)

func main() {
	// Initialize the database object using a configuration file.
	dbConfig, _ := config.ReadDatabaseConfig("./database.json")
	db, err := sql.Open("mysql", dbConfig.DataSourceName())

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Create new router and controllers for handling routing.
	r := httprouter.New()
	uc := controllers.NewUserController(db)
	mc := controllers.NewMasterController(db)

	// Define the routes of the REST service.
	r.POST("/register", uc.CreateUser)
	r.POST("/login", uc.LoginUser)
	r.DELETE("/users", middlewares.AuthMiddleWare(uc.DeleteUser))

	r.GET("/masters", middlewares.AuthMiddleWare(mc.GetMasters))
	r.POST("/masters", middlewares.AuthMiddleWare(mc.CreateMaster))

	r.GET("/masters/:id", middlewares.AuthMiddleWare(mc.GetMaster))
	r.PUT("/masters/:id", middlewares.AuthMiddleWare(mc.UpdateMaster))
	r.DELETE("/masters/:id", middlewares.AuthMiddleWare(mc.DeleteMaster))

	// Start listening to clients.
	http.ListenAndServe("localhost:3000", r)
}
