package usps

import (
	"encoding/json"
	"log"
	"testing"
)

func TestGetAddress(t *testing.T) {

	a := VerificationRequest{
		Address: Address{
			Unit:   "APT 1715",
			Street: "591 Evernia Street",
			City:   "WEST PALM BEACH",
			State:  "FL",
			Zip5:   "33401",
		},
	}

	d, _ := json.Marshal(&a)

	if add, err := GetAddress(string(d)); err != nil {
		t.Error(err)
	} else {
		log.Println(add)
	}

}

func TestGetPostage(t *testing.T) {

	p := RateRequest{
		Packages: []Package{
			{
				Id:             "product-id-1",
				ZipDestination: "90210",
				ZipOrigination: "33401",
				Service:        "PRIORITY",
				Container:      "LG FLAT RATE BOX",
				Pounds:         5,
				Ounces:         10.5,
			},
		},
	}

	d, _ := json.Marshal(p)

	if pp, err := GetPostage(string(d)); err != nil {
		t.Error(err)
	} else {
		log.Println(pp)
	}

}
