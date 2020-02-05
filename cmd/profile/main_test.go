package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"hc-api/pkg/model/profile"
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
	proxy = profile.Proxy{
		token.Value{
			"",
			"",
			[]string{test.CookieValid},
			test.Ip,
		},
		profile.Entity{
			"",
			"connor@wiesow.com",
			"Connor",
			"Van Elswyk",
			"555-555-5555",
		},
	}
)

func TestHandleBadSubject(t *testing.T) {
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleSaveNew(t *testing.T) {
	proxy.Subject = "save"
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleSaveOld(t *testing.T) {
	proxy.Subject = "save"
	proxy.Id = test.ProfileId
	proxy.LastName = "van Elswyk"
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleExpired(t *testing.T) {
	proxy.JwtSlice = []string{test.CookieExpired}
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 402 {
		t.Fail()
	}
}

func TestHandleBadProxyRequest(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{}); out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
