package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"hc-api/pkg/model/offer"
	"hc-api/pkg/model/token"
	"hc-api/test"
	"testing"
)

var (
	requestCtx = events.APIGatewayProxyRequestContext{
		Identity: events.APIGatewayRequestIdentity{
			SourceIP: test.Ip,
		},
	}
	proxy = offer.Proxy{
		token.Value{},
		offer.Entity{
			"",
			"",
			"",
			"1000 lbs, $75/lb, first & final.",
		},
	}
)

func TestHandleInvalidEntity(t *testing.T) {
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleInvalidToken(t *testing.T) {
	proxy.ProductId = test.ProductId
	proxy.JwtSlice = []string{test.CookieExpired}
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 402 {
		t.Fail()
	}
}

func TestHandleSaveNew(t *testing.T) {
	proxy.JwtSlice = []string{test.CookieValid}
	proxy.ProductId = test.ProductId
	proxy.Subject = "save"
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleBadSubject(t *testing.T) {
	proxy.ProductId = test.ProductId
	proxy.Subject = ""
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandle(t *testing.T) {
	go main()
}
