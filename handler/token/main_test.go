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
	if out, _ := Handle(r); out.StatusCode != 200 {
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
		r = events.APIGatewayProxyRequest{Path: "authenticate", QueryStringParameters: map[string]string{"token": out.Body}}
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
		"Id Value",
		0,
		"Issuer Value",
		0,
		"Subject Value",
	}
	b, _ := json.Marshal(&claims)
	r := events.APIGatewayProxyRequest{Path: "authorize", Body: string(b)}
	if out, _ := Handle(r); out.StatusCode == 200 {
		r = events.APIGatewayProxyRequest{Path: "inspect", QueryStringParameters: map[string]string{"token": out.Body}}
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

func TestHandleInspect401(t *testing.T) {
	r := events.APIGatewayProxyRequest{Path: "inspect", QueryStringParameters: map[string]string{"token": "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBdWRpZW5jZSBWYWx1ZSIsImV4cCI6MTU5MTA5MDYzMSwianRpIjoiSWQgVmFsdWUiLCJpYXQiOjE1OTEwOTA2MDcsImlzcyI6Iklzc3VlciBWYWx1ZSIsInN1YiI6IlN1YmplY3QgVmFsdWUifQ.fPchrVG8PIi6txWi9L1VkKOTaHwEfRCwQ1buMLIR_lc; Expires=Tue, 02 Jun 2020 09:37:11 GMT"}}
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
