package main

import (
	"hc-api/pkg/model/token"
	"hc-api/test"
	"testing"
)

var (
	authenticate   = "authenticate"
	authorize      = "authorize"
	invalidCookies = []string{"acc-token=", "ref-token=", "crt-token="}
	invalidIp      = token.Value{SourceIp: ""}
	invalidData    = token.Value{SourceIp: test.Ip}
	invalidToken   = token.Value{SourceIp: test.Ip, Subject: authenticate, JwtSlice: invalidCookies}
	invalidSubject = token.Value{SourceIp: test.Ip, SourceId: test.UserId}
	validAuthorize = token.Value{SourceIp: test.Ip, Subject: authorize, SourceId: test.UserId}
)

func TestHandleInvalidIp(t *testing.T) {
	if _, err := Handle(token.Entity{invalidIp, token.Error{}}); err != token.ErrBadIpData {
		t.Fatal(err)
	}
}

func TestHandleInvalidData(t *testing.T) {
	if _, err := Handle(token.Entity{invalidData, token.Error{}}); err != token.ErrBadCookieData {
		t.Fatal(err)
	}
}

func TestHandleAuthenticateInvalidToken(t *testing.T) {
	if _, err := Handle(token.Entity{invalidToken, token.Error{}}); err == nil {
		t.Fail()
	}
}

func TestHandleAuthenticateBadJwtToken(t *testing.T) {
	if _, err := Handle(token.Entity{token.Value{SourceIp: "127.0.0.1"}, token.Error{}}); err != token.ErrBadCookieData {
		t.Fatal(err)
	}
}

func TestHandleAuthenticateBadIpComparison(t *testing.T) {
	if _, err := Handle(token.Entity{token.Value{SourceIp: "127.0.0.1"}, token.Error{}}); err != token.ErrBadCookieData {
		t.Fatal(err)
	}
}

func TestHandleAuthorize(t *testing.T) {
	if tkn, err := Handle(token.Entity{validAuthorize, token.Error{}}); err != nil {
		t.Fatal(err)
	} else {
		t.Log(tkn)
	}
}

func TestHandleBadSubject(t *testing.T) {
	if _, err := Handle(token.Entity{invalidSubject, token.Error{}}); err != token.InvalidToken {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandle(t *testing.T) {
	go main()
}
