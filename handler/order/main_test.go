package main

import (
	"github.com/aws/aws-lambda-go/events"
	"sam-app/test"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		Path:                  "find",
		Headers:               map[string]string{"Authorize": test.CookieValid},
		QueryStringParameters: map[string]string{"ids": "44ace621-a46e-11ea-8817-2e51bfe26708,bb5a147c-a46b-11ea-b72e-365a3e9d7040"},
	}
	if out, _ := HandleRequest(r); out.StatusCode != 200 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
