package profile

import (
	"hc-api/pkg/model/token"
	"os"
)

type Proxy struct {
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
