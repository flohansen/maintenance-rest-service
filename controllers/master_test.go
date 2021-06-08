package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kluddizz/maintenance-rest-service/config"
	"github.com/kluddizz/maintenance-rest-service/models"
)

var master = models.Master{
  Name: "Test Master",
  Host: "127.0.0.1",
  Port: 5050,
}

var db *sql.DB

func TestMain(m *testing.M) {
  BeforeAll()
  code := m.Run()
  AfterAll()
  os.Exit(code)
}

func BeforeAll() {
  // Read the database config
  dbConfig, err := config.ReadDatabaseConfig("../database.json")

  if err != nil {
    panic(err.Error())
  }

  // Try to open a mysql connection
  db, err = sql.Open("mysql", dbConfig.DataSourceName())

  if err != nil {
    panic(err.Error())
  }

  // Make sure the database is clean before starting tests.
  db.Query("DELETE FROM masters")
}

func AfterAll() {
  // Remove all generated masters from the database.
  db.Query("DELETE FROM masters")

  // Close the database connection.
  db.Close()
}

// Test if the route is able to create new masters.
func TestCreateMaster(t *testing.T) {
  // Create json string from master object
  jsonStr, err := json.Marshal(master)
  if err != nil {
    t.Fatal("Json marshalling failed")
  }

  // Setup new POST request on /masters
  req, err := http.NewRequest("POST", "http://localhost:3000/masters", bytes.NewBuffer([]byte(jsonStr)))
  req.Header.Set("Content-Type", "application/json")

  if err != nil {
    t.Fatal(err.Error())
  }

  // Send the request
  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    t.Fatal(err)
  }

  defer res.Body.Close()

  // Read out the inserted (?) master inside the database for validation later.
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

  // Cleanup the database.
  db.Query("DELETE FROM masters")

  // Check if the status code is the expected one.
  if res.StatusCode != 200 {
    t.Errorf("Expected status code to be %d but received %d", 200, res.StatusCode)
  }

  // Check if the inserted master had the right values.
  if masterDb.Name != master.Name ||
     masterDb.Host != master.Host ||
     masterDb.Port != master.Port {
    t.Errorf("Masters are not equal. Got %+v; want %+v", masterDb, master)
  }
}

// Test if the route is able to delete a valid master from the database and
// returns a success code.
func TestDeleteMaster(t *testing.T) {
  // Insert a new master directly into the dabase.
  queryRes, err := db.Exec(
    "INSERT INTO masters (name, host, port) VALUES (?, ?, ?)",
    master.Name, master.Host, master.Port,
  )

  if err != nil {
    t.Fatalf("Could not insert test master: %s", err.Error())
  }

  // Read out the created id for the master.
  insertedId, _ := queryRes.LastInsertId()

  // Prepare the DELETE request on /masters/:id
  req, err := http.NewRequest(
    "DELETE",
    fmt.Sprintf("http://localhost:3000/masters/%d", insertedId),
    nil,
  )

  if err != nil {
    t.Fatal(err.Error())
  }

  // Send the request
  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    t.Fatal(err.Error())
  }

  defer res.Body.Close()

  // Check wheter the master is still present inside the database.
  var masterExists bool
  err = db.QueryRow(
    "SELECT EXISTS(SELECT 1 FROM masters WHERE id = ?)",
    insertedId,
  ).Scan(&masterExists)

  if err != nil {
    t.Fatal(err.Error())
  }

  // Cleanup the database.
  db.Query("DELETE FROM masters")

  // Check the expected status code.
  if res.StatusCode != 200 {
    t.Errorf("Expected status code to be %d but received %d", 200, res.StatusCode)
  }

  // Check if the master was still present after the DELETE request.
  if masterExists {
    t.Errorf("Expected master existance to be %t but received %t", false, masterExists)
  }
}

// Test if the route returns an error code when invalid master ids are provided.
func TestDeleteMasterFail(t *testing.T) {
  // Prepare the invalid DELETE request on /masters/:id by assuming the id = -1
  // doesn't exist.
  req, err := http.NewRequest(
    "DELETE",
    fmt.Sprintf("http://localhost:3000/masters/%d", -1),
    nil,
  )

  if err != nil {
    t.Fatal(err.Error())
  }

  // Send the invalid request.
  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    t.Fatal(err.Error())
  }

  defer res.Body.Close()

  // Check if the status code is really an error.
  if res.StatusCode != 400 {
    t.Errorf("Expected status code to be %d but received %d", 400, res.StatusCode)
  }
}

func TestUpdateMaster(t *testing.T) {
  // Insert a new master directly into the dabase.
  queryRes, err := db.Exec(
    "INSERT INTO masters (name, host, port) VALUES (?, ?, ?)",
    master.Name, master.Host, master.Port,
  )

  if err != nil {
    t.Fatalf("Could not insert test master: %s", err.Error())
  }

  updatedMaster := master
  updatedMaster.Name = "Updated Master"

  jsonBody, err := json.Marshal(updatedMaster)

  if err != nil {
    t.Fatalf("Could not parse json: %s", err.Error())
  }

  // Read out the created id for the master.
  insertedId, _ := queryRes.LastInsertId()

  // Prepare the PUT request on /masters/:id
  req, err := http.NewRequest(
    "PUT",
    fmt.Sprintf("http://localhost:3000/masters/%d", insertedId),
    bytes.NewBuffer(jsonBody),
  )

  if err != nil {
    t.Fatal(err.Error())
  }

  // Send the request
  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    t.Fatal(err.Error())
  }

  defer res.Body.Close()

  // Read the updated master from the database.
  var updatedMasterDb models.Master
  err = db.QueryRow(
    "SELECT id, name, host, port FROM masters WHERE id = ?",
    insertedId,
  ).Scan(
    &updatedMasterDb.Id, &updatedMasterDb.Name, &updatedMasterDb.Host, &updatedMasterDb.Port,
  )

  if err != nil {
    t.Fatalf("Error while selecting updated master: %s", err.Error())
  }

  // Cleanup the database.
  db.Query("DELETE FROM masters")

  // Check the expected status code.
  if res.StatusCode != 200 {
    t.Errorf("Expected status code to be %d but received %d", 200, res.StatusCode)
  }

  if updatedMasterDb.Name != updatedMaster.Name ||
     updatedMasterDb.Host != master.Host ||
     updatedMasterDb.Port != master.Port {
    t.Errorf("Expected updated master to be %v but received %v", updatedMaster, updatedMasterDb)
  }
}

func TestGetMasters(t *testing.T) {
}

func TestGetSingleMaster(t *testing.T) {
}
