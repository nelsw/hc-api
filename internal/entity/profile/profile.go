package profile

import (
	"hc-api/internal/class"
	"os"
	"strings"
)

type Aggregate struct {
	class.Object
	class.Person
	BrandIds []string `json:"brand_ids"`
}

var table = os.Getenv("PROFILE_TABLE")
var handler = "hcProfileHandler"

func (e *Aggregate) Validate() error {
	return nil
}

func (e *Aggregate) Id() *string {
	s := strings.ToLower(e.Email)
	return &s
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
