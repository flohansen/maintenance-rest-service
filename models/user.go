package models

type (
	User struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Id        string `json:"id"`
	}
)
