package offer

import (
	"hc-api/internal/class"
	"hc-api/internal/entity/product"
	"hc-api/internal/entity/profile"
	"hc-api/internal/entity/token"
	"hc-api/internal/entity/user"
	"os"
	"strings"
)

type Aggregate struct {
	class.Object
	User      user.Aggregate
	Profile   profile.Aggregate
	Product   product.Aggregate
	Token     token.Aggregate
	UserId    string `json:"user_id"`
	ProfileId string `json:"profile_id"`
	Details   string `json:"details"`
}

var table = os.Getenv("OFFER_TABLE")
var handler = "hcOfferHandler"

func (e *Aggregate) Id() *string {
	s := strings.ToLower(e.Profile.Email)
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

func (e *Aggregate) Validate() error {
	return nil
}
