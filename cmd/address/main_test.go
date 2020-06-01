package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"sam-app/pkg/model/address"
	"sam-app/test"
	"testing"
)

var (
	requestCtx = events.APIGatewayProxyRequestContext{
		Identity: events.APIGatewayRequestIdentity{
			SourceIP: test.Ip,
		},
	}
	e = address.Entity{
		test.AddressId,
		"591 Evernia ST",
		"APT 1720",
		"West Palm",
		"FL",
		"33401",
		"",
	}
	proxy = address.Request{}
)

// tests the golden path for the save command
func TestHandleSave200(t *testing.T) {
	proxy.Entity = e
	proxy.Op = "save"
	proxy.Subject = "authenticate"
	proxy.JwtSlice = []string{test.CookieValid}
	b, _ := json.Marshal(&proxy)
	out, err := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if err != nil {
		t.Fatal(err)
	} else if out.StatusCode != 200 {
		t.Fatal(fmt.Errorf("bad response [%v]", out))
	}
}

func TestHandleFindOne200(t *testing.T) {
	proxy.Entity = e
	proxy.Op = "find-one"
	proxy.Subject = "authenticate"
	proxy.JwtSlice = []string{test.CookieValid}
	b, _ := json.Marshal(&proxy)
	out, err := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if err != nil {
		t.Fatal(err)
	} else if out.StatusCode != 200 {
		t.Fatal(fmt.Errorf("bad response [%v]", out))
	}
}

func TestHandleBadOperation(t *testing.T) {
	proxy.Entity = e
	proxy.Op = ""
	proxy.Subject = "authenticate"
	proxy.JwtSlice = []string{test.CookieValid}
	b, _ := json.Marshal(&proxy)
	out, err := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if err != nil {
		t.Fatal(err)
	} else if out.StatusCode != 400 {
		t.Fatal(fmt.Errorf("bad response [%v]", out))
	}
}

// for code coverage purposes only
func TestHandleRequest(t *testing.T) {
	go main()
}
