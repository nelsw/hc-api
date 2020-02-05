package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"hc-api/pkg/model/product"
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
	proxy = product.Proxy{nil, token.Value{}, product.Entity{}}
)

func TestHandleValidateErr(t *testing.T) {
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleTokenErr(t *testing.T) {
	proxy.Name = "foo"
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 402 {
		t.Fail()
	}
}

func TestHandleRequestErr(t *testing.T) {
	proxy.JwtSlice = []string{test.CookieValid}
	proxy.Name = "foo"
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

}

func TestHandleFindOne(t *testing.T) {
	proxy.Subject = "find-one"
	proxy.Id = test.ProductId
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleFindMany(t *testing.T) {
	proxy.Ids = []string{test.ProductId, test.ProductId2}
	proxy.Subject = "find-many"
	b, _ := json.Marshal(&proxy)
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 200 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
