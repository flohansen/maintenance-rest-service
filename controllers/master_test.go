package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

  _ "github.com/go-sql-driver/mysql"
	"github.com/kluddizz/maintenance-rest-service/config"
	"github.com/kluddizz/maintenance-rest-service/models"
)

func TestCreateMaster(t *testing.T) {
  dbConfig, err := config.ReadDatabaseConfig("../database.json")

  if err != nil {
    t.Fatal("Cannot open database config")
  }

  db, err := sql.Open("mysql", dbConfig.DataSourceName())

  if err != nil {
    t.Fatalf("Cannot open database connection: %s", err.Error())
  }

  master := models.Master{
    Name: "Test Master",
    Host: "127.0.0.1",
    Port: 5050,
  }

  jsonStr, err := json.Marshal(master)
  if err != nil {
    t.Fatal("Json marshalling failed")
  }

  req, err := http.NewRequest("POST", "http://localhost:3000/masters", bytes.NewBuffer([]byte(jsonStr)))
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    t.Fatal(err)
  }

  defer res.Body.Close()

  var masterDb models.Master
  err = db.QueryRow(
    "SELECT id, name, host, port FROM masters WHERE name = ?",
    master.Name,
  ).Scan(
    &masterDb.Id, &masterDb.Name, &masterDb.Host, &masterDb.Port,
  )

  if err != nil {
    t.Fatalf("Error in the query: %s", err.Error())
  }

  if masterDb.Name != master.Name ||
     masterDb.Host != master.Host ||
     masterDb.Port != master.Port {
    t.Errorf("Masters are not equal. Got %+v; want %+v", masterDb, master)
  }

}
