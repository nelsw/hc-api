package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"sam-app/pkg/model/order"
	"sam-app/test"
	"testing"
)

func TestHandleFind(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		Path:                  "find",
		Headers:               map[string]string{"Authorize": test.CookieValid},
		QueryStringParameters: map[string]string{"ids": test.OrderIds},
	}
	if out, _ := HandleRequest(r); out.StatusCode != 200 {
		t.Fail()
	}
}

//func TestHandleRates(t *testing.T) {
//	r := events.APIGatewayProxyRequest{
//		Path:                  "rates",
//		Headers:               map[string]string{"Authorize": test.CookieValid},
//		QueryStringParameters: map[string]string{"ids": "44ace621-a46e-11ea-8817-2e51bfe26708,bb5a147c-a46b-11ea-b72e-365a3e9d7040"},
//	}
//	if out, _ := HandleRequest(r); out.StatusCode != 200 {
//		t.Fail()
//	}
//}

func TestHandleSave(t *testing.T) {
	e := order.Entity{
		UserId:    test.UserId,
		AddressId: test.AddressId,
		OrderSum:  0,
		Packages:  nil,
	}
	b, _ := json.Marshal(&e)
	r := events.APIGatewayProxyRequest{
		Path:    "save",
		Headers: map[string]string{"Authorize": test.CookieValid},
		Body:    string(b),
	}
	if out, _ := HandleRequest(r); out.StatusCode != 200 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
