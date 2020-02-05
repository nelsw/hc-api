package product

import (
	"fmt"
	"hc-api/pkg/model/token"
	"os"
)

var (
	ErrCodeBadName = fmt.Errorf("product name must be at least 2 characters in length")
)

type Proxy struct {
	Ids []string `json:"ids"`
	token.Value
	Entity
}

type Entity struct {
	Id          string `json:"id"`
	Sku         string `json:"sku"`
	Img         string `json:"img"`
	Category    string `json:"category"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Other       string `json:"other"`

	AddressId string   `json:"address_id"`
	OwnerId   string   `json:"owner_id"`
	ImageSet  []string `json:"image_set"`
	Quantity  string   `json:"quantity"`
	Stock     string   `json:"stock"`

	// packaging details (calc shipment rates)
	Unit string `json:"unit"` // LB
}

var productTable = os.Getenv("PRODUCT_TABLE")

func (e *Entity) ID() string {
	return e.Id
}

func (*Entity) TableName() string {
	return productTable
}

func (e *Entity) Validate() error {
	if len(e.Name) < 2 {
		return ErrCodeBadName
	}
	return nil
}
