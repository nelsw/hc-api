package ups

import (
	"hc-api/test"
	"log"
	"testing"
)

func TestGetRates(t *testing.T) {
	ty := NodeType{Code: "03", Description: ""}
	p := Package{
		Id: test.ProductId,
		//ZipFrom:  "30067",
		//ZipTo:    "33401",
		//Weight:   7.5,
		//Width:    7,
		//Length:   10,
		//Height:   5,
		NodeType: ty,
		Dimensions: Dimensions{
			Type:   NodeType{Code: "IN"},
			Width:  "7",
			Length: "10",
			Height: "5",
		},
		WeightNode: WeightNode{
			Type:   NodeType{Code: "LBS", Description: "Pounds"},
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

				Weight: WeightNode{
					Type:   NodeType{Code: "LBS", Description: "Pounds"},
					Weight: "7.5",
				},

				Package: p,
			},
		},
	}

	response, err := GetRates(r)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(response)

}
