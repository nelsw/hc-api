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
		QueryStringParameters: map[string]string{
			"token": "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBdWRpZW5jZSBWYWx1ZSIsImV4cCI6MTU5MTE4MTI3MiwianRpIjoiMmY2YTdmM2YtYTQ2OS0xMWVhLWE2MWItNzY3ZTBlODViOTk1IiwiaWF0IjoxNTkxMDk0ODcyLCJpc3MiOiJJc3N1ZXIgVmFsdWUiLCJzdWIiOiJTdWJqZWN0IFZhbHVlIn0.wuFOf1mbNkgnPmz3_iIl-6UFlKw9AkO4IKkvDqFT4Tg; Expires=Wed, 03 Jun 2020 10:47:52 GMT",
		},
		Body: string(b),
	}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleBadRequestNoToken400(t *testing.T) {
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
