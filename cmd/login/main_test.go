package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"sam-app/pkg/model/credential"
	"sam-app/test"
	"testing"
)

func TestHandle400(t *testing.T) {
	if out, err := Handle(events.APIGatewayProxyRequest{}); err != nil || out.StatusCode != 400 {
		t.Fatal(out, err)
	}
}

func TestHandle404(t *testing.T) {
	c := credential.Entity{
		"connor@wiesow.com",
		"",
		"",
		"",
		"Pass123!"}
	b, _ := json.Marshal(&c)
	if out, err := Handle(events.APIGatewayProxyRequest{Body: string(b)}); err != nil || out.StatusCode != 404 {
		t.Fatal(out, err)
	}
}

func TestHandle401(t *testing.T) {
	c := credential.Entity{
		test.CredId,
		test.UserId,
		test.PasswordId,
		"",
		"Pass123!!"}
	b, _ := json.Marshal(&c)
	if out, err := Handle(events.APIGatewayProxyRequest{Body: string(b)}); err != nil || out.StatusCode != 401 {
		t.Fatal(out, err)
	}
}

func TestHandle500(t *testing.T) {
	c := credential.Entity{
		"connor@wiesow.com",
		"",
		"84e5c552-3a40-432a-b5ca-d3ef7a78dee7",
		"",
		"Pass123!"}
	b, _ := json.Marshal(&c)
	if out, err := Handle(events.APIGatewayProxyRequest{Body: string(b)}); err != nil || out.StatusCode != 500 {
		t.Fatal(out, err)
	}
}

func TestHandle200(t *testing.T) {
	c := credential.Entity{
		test.CredId,
		test.UserId,
		test.PasswordId,
		"",
		test.PasswordText}
	b, _ := json.Marshal(&c)
	if out, err := Handle(events.APIGatewayProxyRequest{Body: string(b)}); err != nil || out.StatusCode != 200 {
		t.Fatal(out, err)
	}
}

// for code coverage purposes only
func TestHandle(t *testing.T) {
	go main()
}
