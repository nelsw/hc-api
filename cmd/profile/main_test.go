package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"hc-api/internal"
	"testing"
)

const (
	sourceIp  = "127.0.0.1"
	userId    = "638b13ef-ab84-410a-abb0-c9fd5da45c62"
	profileId = "df3d6542-c828-49e9-a62c-cb6e67d5d730"
)

var (
	session    = internal.NewToken(userId, sourceIp)
	requestCtx = events.APIGatewayProxyRequestContext{Identity: events.APIGatewayRequestIdentity{SourceIP: sourceIp}}
	profile    = Profile{
		Email:     "connor@wiesow.com",
		FirstName: "Connor",
		LastName:  "Van Elswyk",
		Phone:     "555-555-5555",
	}
)

func TestHandleValidSave(t *testing.T) {
	b, _ := json.Marshal(ProfileRequest{Command: "save", Session: session, Profile: profile})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleValidFind(t *testing.T) {
	b, _ := json.Marshal(ProfileRequest{Command: "find", Session: session, Id: profileId})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleBadRequest(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{}); out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadCommand(t *testing.T) {
	b, _ := json.Marshal(ProfileRequest{Command: "", Session: session})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadAuth(t *testing.T) {
	b, _ := json.Marshal(ProfileRequest{Command: "save", Session: "", Profile: profile})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 401 {
		t.Fail()
	}
}

func TestHandleBadSave(t *testing.T) {
	table = ""
	b, _ := json.Marshal(ProfileRequest{Command: "save", Session: session, Profile: profile})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadFind(t *testing.T) {
	table = ""
	b, _ := json.Marshal(ProfileRequest{Command: "find", Session: session, Id: ""})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
