package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"hc-api/internal"
	"testing"
)

const sourceIp = "127.0.0.1"

var (
	session    = internal.NewToken("638b13ef-ab84-410a-abb0-c9fd5da45c62", sourceIp)
	requestCtx = events.APIGatewayProxyRequestContext{Identity: events.APIGatewayRequestIdentity{SourceIP: sourceIp}}
	password   = Password{
		Id:       "0fc6742e-c35b-47a1-a80b-0037bdcd4e8e",
		Password: "Pass123!",
	}
)

func TestHandleOk(t *testing.T) {
	b, _ := json.Marshal(&PasswordRequest{Session: session, Password: password})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 200 {
		t.Fail()
		t.Log(out)
	}
}

func TestHandleBadRequest(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{}); out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadId(t *testing.T) {
	password := Password{Id: "", Password: "Pass123!"}
	b, _ := json.Marshal(&PasswordRequest{Session: session, Password: password})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadSession(t *testing.T) {
	b, _ := json.Marshal(&PasswordRequest{Session: "", Password: password})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 401 {
		t.Fail()
	}
}

func TestHandleBadPassword(t *testing.T) {
	password := Password{Id: password.Id, Password: "Pass1234!"}
	b, _ := json.Marshal(&PasswordRequest{Session: session, Password: password})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 401 {
		t.Fail()
	}
}

func TestHandleBadPasswordLength(t *testing.T) {
	password := Password{Id: password.Id, Password: ""}
	b, _ := json.Marshal(&PasswordRequest{Session: session, Password: password})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadPasswordNumbers(t *testing.T) {
	password := Password{Id: password.Id, Password: "Pass!!!!"}
	b, _ := json.Marshal(&PasswordRequest{Session: session, Password: password})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadPasswordSpecialChar(t *testing.T) {
	password := Password{Id: password.Id, Password: "Pass1234"}
	b, _ := json.Marshal(&PasswordRequest{Session: session, Password: password})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleBadPasswordUppercase(t *testing.T) {
	password := Password{Id: password.Id, Password: "pass123!"}
	b, _ := json.Marshal(&PasswordRequest{Session: session, Password: password})
	out, _ := Handle(events.APIGatewayProxyRequest{RequestContext: requestCtx, Body: string(b)})
	if out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
