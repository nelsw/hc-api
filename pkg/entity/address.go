package entity

import (
	"encoding/json"
	"hc-api/pkg/value"
	"os"
)

type Address struct {
	Id string `json:"id"`
	value.Location
	value.Request
}

var addressHandler = "hcAddressHandler"
var addressTable = os.Getenv("ADDRESS_TABLE")

func (*Address) Table() *string {
	return &addressTable
}

func (*Address) Function() *string {
	return &addressHandler
}

func (e *Address) Ids() []string {
	return []string{e.Id}
}

func (e *Address) Payload() []byte {
	b, _ := json.Marshal(&e)
	return b
}

func (e *Address) Validate() error {
	return nil
}
