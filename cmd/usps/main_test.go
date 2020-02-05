package main

import (
	"hc-api/pkg/model/usps"
	"hc-api/test"
	"testing"
)

func TestHandleValidation(t *testing.T) {
	e := usps.Address{
		"",
		"APT 1715",
		"591 Evernia Street",
		"WEST PALM BEACH",
		"FL",
		"33401",
		"",
	}
	in := usps.Request{"validate", e, []usps.PackageRequest{}}
	if out, err := Handle(in); err != nil {
		t.Fatal(err)
	} else {
		t.Log(out)
	}
}

func TestHandleEstimation(t *testing.T) {
	p := usps.PackageRequest{Id: test.ProductId, ZipTo: "90210", ZipFrom: "33401", Pounds: 5, Ounces: 10.5}
	in := usps.Request{"rate", usps.Address{}, []usps.PackageRequest{p}}
	if out, err := Handle(in); err != nil {
		t.Fatal(err)
	} else {
		t.Log(out)
	}
}

func TestHandleBadRequest(t *testing.T) {
	if out, err := Handle(usps.Request{}); err != ErrOp {
		t.Fatal(out)
	}
}

func TestHandle(t *testing.T) {
	go main()
}
