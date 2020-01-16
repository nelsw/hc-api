package entity

import (
	"encoding/json"
	"os"
)

type User struct {
	Id         string   `json:"id"`
	ProfileId  string   `json:"profile_id"`
	AddressIds []string `json:"address_ids"`
	ProductIds []string `json:"product_ids"`
	OrderIds   []string `json:"order_ids,omitempty"`
	OfferIds   []string `json:"sale_ids,omitempty"`
	SaleIds    []string `json:"sale_ids,omitempty"`
}

var userTable = os.Getenv("USER_TABLE")
var userHandler = "hcUserHandler"

func (e *User) Ids() []string {
	return []string{e.Id}
}

func (e *User) Function() *string {
	return &userHandler
}

func (e *User) Payload() []byte {
	b, _ := json.Marshal(&e)
	return b
}

func (*User) Table() *string {
	return &userTable
}

func (e *User) Validate() error {
	return nil
}
