package entity

import (
	"hc-api/pkg/value"
	"os"
)

type Offer struct {
	value.Authorization
	value.Person
	value.Detail
	UserId    string `json:"user_id"`
	ProfileId string `json:"profile_id"`
}

var offerTable = os.Getenv("OFFER_TABLE")
var offerHandler = "hcOfferHandler"

func (offer *Offer) Ids() []string {
	return nil
}

func (offer *Offer) Payload() []byte {
	return nil
}

func (offer *Offer) Function() *string {
	return &offerHandler
}

func (offer *Offer) Table() *string {
	return &offerTable
}

func (offer *Offer) Validate() error {
	return nil
}
