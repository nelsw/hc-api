package main

import (
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestHandleBadRequest(t *testing.T) {
	r := events.APIGatewayProxyRequest{}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
