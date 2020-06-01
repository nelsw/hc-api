package profile

import (
	"os"
	"sam-app/pkg/model/token"
)

type Request struct {
	Op string `json:"op"`
	token.Value
	Entity
}

type Entity struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

var table = os.Getenv("PROFILE_TABLE")

func (e *Entity) ID() string {
	return e.Id
}

func (*Entity) TableName() string {
	return table
}

func (*Entity) Validate() error {
	return nil
}
