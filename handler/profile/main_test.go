package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"sam-app/pkg/model/profile"
	"sam-app/test"
	"testing"
)

func TestHandleSave200(t *testing.T) {
	e := profile.Entity{
		Id:        test.UserId,
		Email:     "hello@gmail.com",
		FirstName: "Jimmy",
		LastName:  "Kowalski",
		Phone:     "555-555-5555",
	}
	b, _ := json.Marshal(&e)
	r := events.APIGatewayProxyRequest{
		Path: "save",
		Headers: map[string]string{
			"Authorize": test.CookieValid,
		},
		Body: string(b),
	}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleBadNoToken401(t *testing.T) {
	e := profile.Entity{
		Id:        test.UserId,
		Email:     "hello@gmail.com",
		FirstName: "Jimmy",
		LastName:  "Kowalski",
		Phone:     "555-555-5555",
	}
	b, _ := json.Marshal(&e)
	r := events.APIGatewayProxyRequest{
		Path: "save",
		Headers: map[string]string{
			"Authorize": "foo",
		},
		Body: string(b),
	}
	if out, _ := Handle(r); out.StatusCode != 401 {
		t.Fail()
	}
}

func TestHandleBadPath(t *testing.T) {
	e := profile.Entity{
		Id:        test.UserId,
		Email:     "hello@gmail.com",
		FirstName: "Jimmy",
		LastName:  "Kowalski",
		Phone:     "555-555-5555",
	}
	b, _ := json.Marshal(&e)
	r := events.APIGatewayProxyRequest{
		Path: "foo",
		Headers: map[string]string{
			"Authorize": test.CookieValid,
		},
		Body: string(b),
	}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadRequest400(t *testing.T) {
	r := events.APIGatewayProxyRequest{}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandle(t *testing.T) {
	go main()
}
