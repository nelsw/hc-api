package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
	"sam-app/test"
	"testing"
)

func TestHandleAuthorize200(t *testing.T) {
	claims := jwt.StandardClaims{
		"Audience Value",
		0,
		test.UserId,
		0,
		"Issuer Value",
		0,
		"Subject Value",
	}
	b, _ := json.Marshal(&claims)
	r := events.APIGatewayProxyRequest{Path: "authorize", Body: string(b)}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleAuthenticate200(t *testing.T) {
	claims := jwt.StandardClaims{
		"Audience Value",
		0,
		test.UserId,
		0,
		"Issuer Value",
		0,
		"Subject Value",
	}
	b, _ := json.Marshal(&claims)
	r := events.APIGatewayProxyRequest{Path: "authorize", Body: string(b)}
	if out, _ := Handle(r); out.StatusCode == 200 {
		r = events.APIGatewayProxyRequest{Path: "authenticate", Headers: map[string]string{"Authorize": out.Body}}
		if out, _ := Handle(r); out.StatusCode == 200 {
			return
		}
	}
	t.Fail()
}

func TestHandleInspect200(t *testing.T) {
	claims := jwt.StandardClaims{
		"Audience Value",
		0,
		test.UserId,
		0,
		"Issuer Value",
		0,
		"Subject Value",
	}
	b, _ := json.Marshal(&claims)
	r := events.APIGatewayProxyRequest{Path: "authorize", Body: string(b)}
	if out, _ := Handle(r); out.StatusCode == 200 {
		r = events.APIGatewayProxyRequest{Path: "inspect", Headers: map[string]string{"Authorize": out.Body}}
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
	} else {
		t.Log(out)
	}
}

func TestHandleAuthenticate401(t *testing.T) {
	r := events.APIGatewayProxyRequest{Path: "authenticate", Headers: map[string]string{"Authorize": "token=foo"}}
	if out, _ := Handle(r); out.StatusCode != 401 {
		t.Fail()
	} else {
		t.Log(out)
	}
}

func TestHandleInspect401(t *testing.T) {
	r := events.APIGatewayProxyRequest{Path: "inspect", Headers: map[string]string{"Authorize": test.CookieExpired}}
	if out, _ := Handle(r); out.StatusCode != 401 {
		t.Fail()
	} else {
		t.Log(out)
	}
}

func TestHandleBadRequest400(t *testing.T) {
	r := events.APIGatewayProxyRequest{}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	} else {
		t.Log(out)
	}
}

// for code coverage purposes only
func TestHandle(t *testing.T) {
	go main()
}
