package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/model/profile"
	"testing"
)

var body, tkn string

func init() {
	data, _ := json.Marshal(&profile.Entity{
		"testUserId",
		"hello@gmail.com",
		"Jimmy",
		"Kowalski",
		"555-555-5555",
		"https://www.celebritynooz.com/img/barrynewman-then.jpg",
	})
	body = string(data)

	claims, _ := json.Marshal(&jwt.StandardClaims{
		"127.0.0.1",
		0,
		"testUserId",
		0,
		"testForProfile",
		0,
		"test",
	})
	auth := client.Invoke("tokenHandler", events.APIGatewayProxyRequest{Path: "authorize", Body: string(claims)})
	tkn = auth.Body
}

func TestHandleSave200(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		Headers:               map[string]string{"Authorize": tkn},
		QueryStringParameters: map[string]string{"path": "save"},
		Body:                  body,
	}); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleBadNoToken401(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		Headers:               map[string]string{"Authorize": "foo"},
		QueryStringParameters: map[string]string{"path": "save"},
		Body:                  body,
	}); out.StatusCode != 401 {
		t.Fail()
	}
}

func TestHandleBadPath(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		Headers: map[string]string{"Authorize": tkn},
		Body:    body,
	}); out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadRequest400(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{}); out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandle(t *testing.T) {
	go main()
}
