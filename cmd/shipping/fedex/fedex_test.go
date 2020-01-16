package fedex

import (
	"encoding/json"
	"log"
	"testing"
)

// test url https://wsbeta.fedex.com:443/web-services
// test office  integrator id   123
// test client  product id      TEST
// test client  product version 9999
func TestGenerateSOAPRequest(t *testing.T) {

	r := PackageRequest{
		Id:                 "product-id-1",
		ShipperZipCode:     "33401",
		ShipperStateCode:   "FL",
		RecipientZipCode:   "90210",
		RecipientStateCode: "CA",
		Weight:             7.5,
		Length:             10,
		Width:              7,
		Height:             5,
	}

	var rr RateRequest
	rr.Packages = append(rr.Packages, r)
	b, err := json.Marshal(&rr)
	if err != nil {
		t.Error(err)
	}
	out, err := GetRates(string(b))
	if err != nil {
		t.Error(err)
	} else {
		log.Println(out)
	}

}
