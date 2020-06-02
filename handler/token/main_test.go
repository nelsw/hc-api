package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

func TestHandleAuthorize200(t *testing.T) {
	claims := jwt.StandardClaims{
		"Audience Value",
		0,
		"Id Value",
		0,
		"Issuer Value",
		0,
		"Subject Value",
	}
	b, _ := json.Marshal(&claims)
	r := events.APIGatewayProxyRequest{Path: "authorize", Body: string(b)}
	if out, _ := Handle(r); out.StatusCode == 200 {
		t.Log(out.Body)
	} else {
		t.Fail()
	}
}

func TestHandleAuthenticate200(t *testing.T) {
	claims := jwt.StandardClaims{
		"Audience Value",
		0,
		"Id Value",
		0,
		"Issuer Value",
		0,
		"Subject Value",
	}
	b, _ := json.Marshal(&claims)
	r := events.APIGatewayProxyRequest{Path: "authorize", Body: string(b)}
	if out, _ := Handle(r); out.StatusCode == 200 {
		token := out.Body
		t.Log(token)
		r = events.APIGatewayProxyRequest{Path: "authenticate", QueryStringParameters: map[string]string{"token": token}}
		if out, _ := Handle(r); out.StatusCode == 200 {
			t.Log(out.Body)
			return
		}
	}
	t.Fail()
}

func TestHandleAuthorize400(t *testing.T) {
	r := events.APIGatewayProxyRequest{Path: "authorize"}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleAuthenticate401(t *testing.T) {
	r := events.APIGatewayProxyRequest{Path: "authenticate", QueryStringParameters: map[string]string{"token": "token=foo"}}
	if out, _ := Handle(r); out.StatusCode != 401 {
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
