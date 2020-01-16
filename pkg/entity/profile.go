package entity

import (
	"encoding/json"
	"hc-api/pkg/value"
	"os"
)

type Profile struct {
	value.Authorization
	value.Person
	value.Request
	Id       string   `json:"id"`
	BrandIds []string `json:"brand_ids"`
}

var profileTable = os.Getenv("PROFILE_TABLE")
var profileHandler = "hcProfileHandler"

func (e *Profile) Ids() []string {
	return []string{e.Id}
}

func (e *Profile) Validate() error {
	return nil
}

func (e *Profile) Payload() []byte {
	b, _ := json.Marshal(&e)
	return b
}

func (e *Profile) Function() *string {
	return &profileHandler
}

func (*Profile) Table() *string {
	return &profileTable
}
