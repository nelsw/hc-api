package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"sam-app/pkg/model/credential"
	"testing"
)

func TestHandleLogin200(t *testing.T) {
	credentials := credential.Entity{"hello@gmail.com", "Pass123!"}
	b, _ := json.Marshal(&credentials)
	r := events.APIGatewayProxyRequest{Body: string(b)}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	} else {
		t.Log(out)
	}
}

func TestHandleLogin401(t *testing.T) {
	credentials := credential.Entity{"hello@gmail.com", "Pass1234!"}
	b, _ := json.Marshal(&credentials)
	r := events.APIGatewayProxyRequest{Body: string(b)}
	if out, _ := Handle(r); out.StatusCode != 401 {
		t.Fail()
	} else {
		t.Log(out)
	}
}

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
