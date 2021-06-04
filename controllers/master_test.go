package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kluddizz/maintenance-rest-service/config"
	"github.com/kluddizz/maintenance-rest-service/models"
)

func TestCreateMaster(t *testing.T) {
  dbConfig, _ := config.ReadDatabaseConfig("./database.json")
  db, err := sql.Open("mysql", dbConfig.DataSourceName())

  if err != nil {
    t.Fatal("Cannot open database connection")
  }

  user := models.User{
    UserName: "TestUser",
    FirstName: "Test",
    LastName: "User",
    Email: "test.user@mail.com",
  }

  jsonStr, err := json.Marshal(user)
  if err != nil {
    t.Fatal("Json marshalling failed")
  }

  req, err := http.NewRequest("POST", "http://localhost:3000", bytes.NewBuffer([]byte(jsonStr)))
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  res, err := client.Do(req)

  if err != nil {
    t.Fatal(err)
  }

  defer res.Body.Close()

  var userDb models.UserDb
  err = db.QueryRow(
    "SELECT id, username, firstname, lastname, email FROM users WHERE username = ?",
    user.UserName,
  ).Scan(
    &userDb.Id, &userDb.UserName, &userDb.FirstName, &userDb.LastName, &userDb.Email,
  )

  if err != nil {
    t.Fatal("Error in the query")
  }

  if userDb.Id != user.Id ||
    userDb.UserName != user.UserName ||
    userDb.FirstName != user.FirstName ||
    userDb.LastName != user.LastName ||
    userDb.Email != user.Email {
    t.Errorf("Users are not equal. Got %+v; want %+v", userDb.UserName, user.UserName)
  }

}
