package user

import (
	"hc-api/internal/class"
	"os"
)

type Aggregate struct {
	class.Object
	User
}

// Primary user object for the domain, visible to client and server.
// Each property is a reference to a unique ID, or collection of unique ID's.
type User struct {
	ProfileId  string   `json:"profile_id"`
	AddressIds []string `json:"address_ids"`
	ProductIds []string `json:"product_ids"`
	OrderIds   []string `json:"order_ids,omitempty"`
	OfferIds   []string `json:"sale_ids,omitempty"`
	SaleIds    []string `json:"sale_ids,omitempty"`
}

var table = os.Getenv("USER_TABLE")
var handler = "hcUserHandler"

func (e *Aggregate) Id() *string {
	return e.UUID()
}

func (e *Aggregate) Payload() []byte {
	return nil
}

func (*Aggregate) Handler() *string {
	return &handler
}

func (*Aggregate) Name() *string {
	return &table
}

func (e *Aggregate) Validate() error {
	return nil
}
