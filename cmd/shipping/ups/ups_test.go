package ups

import (
	"encoding/json"
	"log"
	"testing"
)

func TestGetRates(t *testing.T) {
	ty := Type{Code: "03", Description: ""}
	p := Package{
		Id:             "product-id-1",
		ZipOrigination: "30067",
		ZipDestination: "33401",
		Weight:         7.5,
		Width:          7,
		Length:         10,
		Height:         5,
		Type:           ty,
		Dimensions: Dimensions{
			Type:   Type{Code: "IN"},
			Width:  "7",
			Length: "10",
			Height: "5",
		},
		PackageWeight: PackageWeight{
			Type:   Type{Code: "LBS", Description: "Pounds"},
			Weight: "7",
		},
	}
	r := PostageRateRequest{
		Packages: []Package{p},
		//order
		RateRequest: RateRequest{
			//package
			Shipment: Shipment{
				//ShippingEntity
				Shipper: ShippingEntity{
					Address: Address{
						PostalCode:  "30067",
						CountryCode: "US",
					},
				},
				//ShippingEntity
				ShipTo: ShippingEntity{
					Address: Address{
						PostalCode:  "95113",
						CountryCode: "US",
					},
				},
				//ShippingEntity
				ShipFrom: ShippingEntity{
					Address: Address{
						PostalCode:  "30067",
						CountryCode: "US",
					},
				},

				Service: ty,

				Weight: ShipmentTotalWeight{
					Type:   Type{Code: "LBS", Description: "Pounds"},
					Weight: "7.5",
				},

				Package: p,
			},
		},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	response, err := GetRates(string(b))
	if err != nil {
		t.Fatal(err)
	}
	log.Println(response)

}
