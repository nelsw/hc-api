package main

import (
	"sam-app/pkg/model/ups"
	"sam-app/test"
	"testing"
)

func TestHandleRequest(t *testing.T) {

	lbs := ups.NodeType{"LBS", "Pounds"}
	inc := ups.NodeType{"IN", "Inches"}

	svc := ups.NodeType{"03", "GROUND"}
	w := ups.WeightNode{lbs, "7.5"}
	d := ups.Dimensions{inc, "7", "10", "5"}

	p := ups.Package{test.ProductId, svc, d, w}

	fr := ups.ShippingEntity{ups.Address{"90210", "US"}}
	to := ups.ShippingEntity{ups.Address{"33401", "US"}}

	s := ups.Shipment{fr, to, fr, svc, w, p}

	r := ups.PostageRateRequest{ups.RateRequest{s}, []ups.Package{p}}

	if out, err := Handle(r); err != nil {
		t.Fatal(err)
	} else {
		t.Log(out)
	}
}

func TestHandle(t *testing.T) {
	go main()
}
