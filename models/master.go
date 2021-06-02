package models

type (
	Master struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Host string `json:"host"`
		Port int    `json:"port"`
	}
)
