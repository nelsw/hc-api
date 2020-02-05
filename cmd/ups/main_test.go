package main

import (
	"hc-api/pkg/model/ups"
	"hc-api/test"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	ty := ups.NodeType{Code: "03", Description: ""}
	p := ups.Package{
		Id:       test.ProductId,
		ZipFrom:  "90210",
		ZipTo:    "33401",
		Weight:   7.5,
		Width:    7,
		Length:   10,
		Height:   5,
		NodeType: ty,
		Dimensions: ups.Dimensions{
			Type:   ups.NodeType{Code: "IN"},
			Width:  "7",
			Length: "10",
			Height: "5",
		},
		WeightNode: ups.WeightNode{
			Type:   ups.NodeType{Code: "LBS", Description: "Pounds"},
			Weight: "7",
		},
	}
	r := ups.PostageRateRequest{
		Packages: []ups.Package{p},
		//order
		RateRequest: ups.RateRequest{
			//package
			Shipment: ups.Shipment{
				//ShippingEntity
				Shipper: ups.ShippingEntity{
					Address: ups.Address{
						PostalCode:  "90210",
						CountryCode: "US",
					},
				},
				//ShippingEntity
				ShipTo: ups.ShippingEntity{
					Address: ups.Address{
						PostalCode:  "33401",
						CountryCode: "US",
					},
				},
				//ShippingEntity
				ShipFrom: ups.ShippingEntity{
					Address: ups.Address{
						PostalCode:  "90210",
						CountryCode: "US",
					},
				},

				Service: ty,

				Weight: ups.WeightNode{
					Type:   ups.NodeType{Code: "LBS", Description: "Pounds"},
					Weight: "7.5",
				},

				Package: p,
			},
		},
	}

	if out, err := Handle(r); err != nil {
		t.Fatal(err)
	} else {
		t.Log(out)
	}
}

func TestHandle(t *testing.T) {
	go main()
}
