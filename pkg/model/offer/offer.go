package offer

import (
	"fmt"
	"hc-api/pkg/model/token"
	"os"
)

var ErrCodeDetailEmpty = fmt.Errorf("offer detail empty")
var ErrCodeProductIdEmpty = fmt.Errorf("product id empty")

type Proxy struct {
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
