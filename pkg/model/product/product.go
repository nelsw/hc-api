package product

import (
	"os"
)

type Entity struct {
	Id          string   `json:"id"`
	Category    string   `json:"category"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	Images      []string `json:"images"`
	OwnerId     string   `json:"owner_id"`
	AddressId   string   `json:"address_id"` // shipping departure location
	Unit        string   `json:"unit"`       // LB, OZ, etc.
	Weight      int64    `json:"weight"`
	Stock       int      `json:"stock"`
}

var table = os.Getenv("PRODUCT_TABLE")

func (e *Entity) ID() string {
	return e.Id
}

func (*Entity) TableName() string {
	return table
}

func (e *Entity) Validate() error {
	return nil
}
