package address

import (
	"encoding/json"
	"hc-api/internal/class"
	"os"
)

type Aggregate struct {
	class.Token
	class.Object
	class.Address
	class.Request
}

var handler = "hcAddressHandler"
var table = os.Getenv("ADDRESS_TABLE")

func (*Aggregate) Name() *string {
	return &table
}

func (*Aggregate) Handler() *string {
	return &handler
}

func (e *Aggregate) Id() *string {
	return e.UUID()
}

func (e *Aggregate) Payload() []byte {
	b, _ := json.Marshal(&e)
	return b
}

func (e *Aggregate) Validate() error {
	return nil
}
