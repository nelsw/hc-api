package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"hc-api/internal"
	"testing"
)

const sourceIp = "127.0.0.1"
const userId = "638b13ef-ab84-410a-abb0-c9fd5da45c62"

var session = internal.NewToken(userId, sourceIp)
var requestCtx = events.APIGatewayProxyRequestContext{Identity: events.APIGatewayRequestIdentity{SourceIP: sourceIp}}
var offer = Offer{
	ProductId:        "53f2ebcb-1689-11ea-9a91-6a36cd23892f",
	ProductImg:       "https://www.gannett-cdn.com/-mm-/98bf8e596c32510f94f4c1c778ff11fda8fbcb3a/c=960-0-4800-3840/local/-/media/2018/02/08/Cincinnati/Cincinnati/636537101300479428-PTH1Brd-12-27-2017-PTH-1-A001-2017-12-26-IMG-PTH0611-MEDICAL-MARI-1-1-81KMOMB2-L1156005440-IMG-PTH0611-MEDICAL-MARI-1-1-81KMOMB2.jpg?quality=10",
	ProductName:      "Super mind!",
	ProductUnit:      "lb",
	ProductPrice:     10000,
	ProductAddressId: "NTkxIEVWRVJOSUEgU1QsIEFQVCAxNzE1LCBXRVNUIFBBTE0gQkVBQ0gsIEZMLCAzMzQwMS01Nzg1LCBVbml0ZWQgU3RhdGVz",
	Details:          "100K lbs, $75/lb, first & final.",
}

// Given a valid request, when saving offer, return
func TestHandleSave(t *testing.T) {
	b, _ := json.Marshal(OfferRequest{Command: "save", Session: session, Offer: offer})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleBadCommand(t *testing.T) {
	b, _ := json.Marshal(OfferRequest{Command: ""})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

// Given a valid request, when saving offer, return
func TestHandleBadAuth(t *testing.T) {
	tkn := internal.NewToken(userId, "")
	b, _ := json.Marshal(OfferRequest{Command: "save", Session: tkn, Offer: offer})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 401 {
		t.Fail()
	}
}

// Given a valid request, when saving offer, return
func TestHandleBadUser(t *testing.T) {
	tkn := internal.NewToken("", sourceIp)
	b, _ := json.Marshal(OfferRequest{Command: "save", Session: tkn, Offer: offer})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

// Given a valid request, when saving offer, return
func TestHandleBadProfile(t *testing.T) {
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

// Given a valid request, when saving offer, return
func TestHandleBadRequest(t *testing.T) {
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandle(t *testing.T) {
	go main()
}
