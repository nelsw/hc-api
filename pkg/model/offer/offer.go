package offer

import (
	"fmt"
	"os"
	"sam-app/pkg/model/token"
)

var ErrCodeDetailEmpty = fmt.Errorf("offer detail empty")
var ErrCodeProductIdEmpty = fmt.Errorf("product id empty")

type Request struct {
	Op string `json:"op"`
	token.Value
	Entity
}

type Entity struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Detail    string `json:"detail"`
}

var table = os.Getenv("OFFER_TABLE")

func (e *Entity) ID() string {
	return e.Id
}

func (e *Entity) Validate() error {
	if e.ProductId == "" {
		return ErrCodeProductIdEmpty
	} else if e.Detail == "" {
		return ErrCodeDetailEmpty
	}
	return nil
}

func (e *Entity) TableName() string {
	return table
}
