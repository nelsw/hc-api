package order

import (
	"os"
	"sam-app/pkg/model/token"
)

type Request struct {
	Op  string   `json:"op"`
	Ids []string `json:"ids"`
	token.Value
	Entity
}

type Entity struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AddressId string    `json:"address_id"`
	OrderSum  int64     `json:"order_sum"`
	Packages  []Package `json:"packages"`
}

type Package struct {
	Id string `json:"id"` // product id

	ProductPrice int64 `json:"product_price"`
	ProductQty   int   `json:"product_qty"`
	ProductSum   int64 `json:"product_sum"`

	AddressId string `json:"address_id"`

	ZipOrigination string `json:"zip_origination"` // usps
	ZipDestination string `json:"zip_destination"` // usps

	RecipientStateCode string `json:"recipient_state_code"` // fedex
	ShipperStateCode   string `json:"shipper_state_code"`   // fedex

	// ups, fedex, usps
	ProductPounds int     `json:"pounds"`
	ProductOunces float32 `json:"ounces"`
	ProductWeight float32 `json:"product_weight"`
	ProductLength int     `json:"product_length"`
	ProductWidth  int     `json:"product_width"`
	ProductHeight int     `json:"product_height"`

	// ups, fedex
	TotalLength int     `json:"length"`
	TotalWidth  int     `json:"width"`
	TotalHeight int     `json:"height"`
	TotalWeight float32 `json:"weight"`

	ShipVendor  string `json:"ship_name"`
	ShipService string `json:"ship_service"`
	ShipRate    int64  `json:"ship_rate"`

	TotalPrice int64 `json:"total_price"`
}

var table = os.Getenv("ORDER_TABLE")

func (e Entity) ID() string {
	return e.Id
}

func (e *Entity) TableName() string {
	return table
}

func (e *Entity) Validate() error {
	return nil
}
