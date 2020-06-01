package main

import (
	"sam-app/pkg/model/fedex"
	"sam-app/test"
	"testing"
)

func TestGenerateSOAPRequest(t *testing.T) {

	r := fedex.Package{
		Id:               test.ProductId,
		ShipperZipCode:   "33401",
		ShipperStateCode: "FL",
		Weight:           7.5,
		Length:           10,
		Width:            7,
		Height:           5,
	}

	in := fedex.RateRequest{"CA", "90210", []fedex.Package{r}}

	if out, err := Handle(in); err != nil {
		t.Fatal(err)
	} else {
		t.Log(out)
	}

}

func TestHandle(t *testing.T) {
	go main()
}
