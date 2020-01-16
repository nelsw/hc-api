package entity

import (
	"encoding/json"
	"fmt"
	"hc-api/pkg/value"
	"os"
)

type Product struct {
	value.Authorization
	value.Request
	value.Detail
	Id        string   `json:"id"`
	AddressId string   `json:"address_id"`
	OwnerId   string   `json:"owner_id"`
	ImageSet  []string `json:"image_set"`
	Quantity  string   `json:"quantity"`
	Stock     string   `json:"stock"`
	Deleted   string   `json:"deleted,omitempty"`
	// packaging details (calc shipping rates)
	Unit   string  `json:"unit"` // LB
	Weight float32 `json:"weight"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Length int     `json:"length"`
}

var productTable = os.Getenv("PRODUCT_TABLE")
var productHandler = "hcProductHandler"

func (e *Product) Ids() []string {
	return []string{e.Id}
}

func (e *Product) Validate() error {
	switch e.Case {
	case "find-by-id":
	case "find-by-ids":
		if len(e.Ids()) == 0 {
			return fmt.Errorf("empty id\n")
		}
	case "find-by-brand-id":
		if e.BrandId == "" {
			return fmt.Errorf("empty brand id\n")
		}
	case "update":
		if e.Expression == "" {
			return fmt.Errorf("bad expression=[%s]", e.Expression)
		}
	case "save":
		if len(e.Name) < 3 {
			return fmt.Errorf("bad name [%s], must be > 2 characters in length", e.Name)
		} else if e.Price < 0 {
			return fmt.Errorf("bad price [%d]", e.Price)
		}
	}
	return nil
}

func (e *Product) Payload() []byte {
	if e.Request.Case == "save" {
		b, _ := json.Marshal(&e)
		return b
	} else {
		b, _ := json.Marshal(&e.Request)
		return b
	}
}

func (*Product) Function() *string {
	return &productHandler
}

func (*Product) Table() *string {
	return &productTable
}
