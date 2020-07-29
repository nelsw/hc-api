package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/model/address"
	"sam-app/pkg/model/usps"
	"testing"
)

var addressBody, packageBody, tkn string

func init() {
	addressData, _ := json.Marshal(&address.Entity{
		"",
		"APT 1715",
		"591 EVERNIA ST",
		"WEST PALM",
		"FL",
		"33401",
		"",
	})
	addressBody = string(addressData)

	packageData, _ := json.Marshal([]usps.PackageRequest{
		{
			Id:      "testProductId",
			ZipTo:   "34210",
			ZipFrom: "33401",
			Pounds:  5,
			Ounces:  10.5,
		},
	})
	packageBody = string(packageData)

	claims, _ := json.Marshal(&jwt.StandardClaims{
		"127.0.0.1",
		0,
		"testUserId",
		0,
		"usps/main_test.go",
		0,
		"init",
	})
	auth := client.Invoke("tokenHandler", events.APIGatewayProxyRequest{Path: "authorize", Body: string(claims)})
	tkn = auth.Body
}

func TestHandleValidation(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		Headers:               map[string]string{"Authorize": tkn},
		QueryStringParameters: map[string]string{"path": "validate"},
		Body:                  addressBody,
	}); out.StatusCode != 200 {
		t.Fatal(out)
	} else {
		t.Log(out)
	}
}

func TestHandleUnauthorized(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		Headers:               map[string]string{"Authorize": "token=foo"},
		QueryStringParameters: map[string]string{"path": "validate"},
		Body:                  addressBody,
	}); out.StatusCode != 401 {
		t.Fatal(out)
	} else {
		t.Log(out)
	}
}

func TestHandleEstimation(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		Headers:               map[string]string{"Authorize": tkn},
		QueryStringParameters: map[string]string{"path": "rate"},
		Body:                  packageBody,
	}); out.StatusCode != 200 {
		t.Fatal(out)
	} else {
		t.Log(out)
	}
}

func TestHandleBadRequest(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		Headers:               map[string]string{"Authorize": tkn},
		QueryStringParameters: map[string]string{"path": ""},
		Body:                  packageBody,
	}); out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandle(t *testing.T) {
	go main()
}
