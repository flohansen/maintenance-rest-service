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
	"github.com/kluddizz/maintenance-rest-service/utils"
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

func TestUpdateMasterFail(t *testing.T) {
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
    fmt.Sprintf("http://localhost:3000/masters/%d", -1),
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
  if res.StatusCode != 400 {
    t.Errorf("Expected status code to be %d but received %d", 400, res.StatusCode)
  }

  if updatedMasterDb.Name != master.Name ||
     updatedMasterDb.Host != master.Host ||
     updatedMasterDb.Port != master.Port {
    t.Errorf("Expected updated master to be %v but received %v", updatedMaster, updatedMasterDb)
  }
}

// Test if we get all the masters stored in the database.
func TestGetMasters(t *testing.T) {
  numberMasters := 5

  // Insert some masters into the database.
  for i := 0; i < numberMasters; i++ {
    db.Exec(
      "INSERT INTO masters (name, host, port) VALUES (?, ?, ?)",
      fmt.Sprintf("Test Master %d", i), master.Host, master.Port,
    )
  }

  // Prepare request.
  req, err := http.NewRequest(
    "GET",
    "http://localhost:3000/masters",
    nil,
  )

  if err != nil {
    t.Fatalf("Cannot create HTTP request: %s", err.Error())
  }

  // Send request.
  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    t.Fatalf("Error while sending request: %s", err.Error())
  }

  defer res.Body.Close()
  
  // Parse response.
  var response models.JsonResponse
  decoder := json.NewDecoder(res.Body)
  err = decoder.Decode(&response)

  if err != nil {
    t.Fatalf("Could not parse json: %s", err.Error())
  }

  // Cleanup database.
  db.Exec("DELETE FROM masters")

  // Check if request was successful.
  if res.StatusCode != 200 {
    t.Errorf("Expected status code to be %d but received %d", 200, res.StatusCode)
  }

  // Check if the amount of fetched masters.
  masters := response.Content.([]interface{})
  if len(masters) != numberMasters {
    t.Errorf("Expected number of masters to be %d but received %d", numberMasters, len(masters))
  }
}

// Test if we can fetch masters by id.
func TestGetSingleMaster(t *testing.T) {
  // Insert one master into the database.
  queryRes, _ := db.Exec(
    "INSERT INTO masters (name, host, port) VALUES (?, ?, ?)",
    master.Name, master.Host, master.Port,
  )

  // Get the auto generated master id.
  insertedId, _ := queryRes.LastInsertId()

  // Create new GET request.
  req, _ := http.NewRequest(
    "GET",
    fmt.Sprintf("http://localhost:3000/masters/%d", insertedId),
    nil,
  )

  // Send request.
  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    t.Fatalf("Cannot send request: %s", err.Error())
  }

  defer res.Body.Close()

  // Parse the response.
  var response models.JsonResponse
  decoder := json.NewDecoder(res.Body)
  decoder.Decode(&response)

  var masterResponse models.Master
  utils.MapToStruct(response.Content, &masterResponse)

  // Cleanup database.
  db.Exec("DELETE FROM masters")

  // Check status code.
  if res.StatusCode != 200 {
    t.Errorf("Expected status code to be %d but received %d", 200, res.StatusCode)
  }

  // Check fetched master.
  if masterResponse.Host != master.Host ||
     masterResponse.Name != master.Name ||
     masterResponse.Port != master.Port {
    t.Errorf("Expected master to be %v but received %v", master, masterResponse)
  }
}

// Test if we cannot fetch masters by invalid ids.
func TestGetSingleMasterFail(t *testing.T) {
  // Create new GET request.
  req, _ := http.NewRequest(
    "GET",
    fmt.Sprintf("http://localhost:3000/masters/%d", -1),
    nil,
  )

  // Send request.
  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    t.Fatalf("Cannot send request: %s", err.Error())
  }

  defer res.Body.Close()

  // Check status code.
  if res.StatusCode != 400 {
    t.Errorf("Expected status code to be %d but received %d", 400, res.StatusCode)
  }
}

