package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type (
	DatabaseConfig struct {
		Database string `json:"database"`
		UserName string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	}
)

func ReadDatabaseConfig(path string) (*DatabaseConfig, error) {
	dbConfig := DatabaseConfig{}

	jsonString, err := ioutil.ReadFile(path)
	if err != nil {
		return &dbConfig, err
	}

	err = json.Unmarshal(jsonString, &dbConfig)
	if err != nil {
		return &dbConfig, err
	}

	return &dbConfig, nil
}

func (dbConfig *DatabaseConfig) DataSourceName() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		dbConfig.UserName,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)
}
