package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	addressId = "NTkxIEVWRVJOSUEgU1QsIEFQVCAxNzA1LCBXRVNUIFBBTE0gQkVBQ0gsIEZMLCAzMzQwMS01Nzg0LCBVbml0ZWQgU3RhdGVz"
	userId    = "638b13ef-ab84-410a-abb0-c9fd5da45c62"
	sourceIp  = "127.0.0.1"
)

var (
	jwtKey  = []byte(os.Getenv("JWT_KEY"))
	address = Address{
		Id:     addressId,
		Street: "591 Evernia ST",
		Unit:   "APT 1715",
		City:   "West Palm",
		State:  "FL",
		Zip5:   "33401",
		Zip4:   "5784",
	}
	addressRequest = AddressRequest{Ids: []string{addressId}, Address: address}
	requestCtx     = events.APIGatewayProxyRequestContext{Identity: events.APIGatewayRequestIdentity{SourceIP: sourceIp}}
)

func init() {
	type Claims struct {
		Id string `json:"id"`
		Ip string `json:"ip"`
		jwt.StandardClaims
	}
	expiry := time.Now().Add(30 * time.Minute)
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{Id: userId, Ip: sourceIp, StandardClaims: jwt.StandardClaims{ExpiresAt: expiry.Unix()}})
	str, _ := tkn.SignedString(jwtKey)
	cookie := &http.Cookie{Name: "token", Value: str, Expires: expiry, HttpOnly: false}
	addressRequest.AccessToken = cookie.String()
}

// tests the golden path for the save command
func TestHandleRequestSave200(t *testing.T) {
	addressRequest.Command = "save"
	requestBytes, _ := json.Marshal(&addressRequest)
	if r, err := HandleRequest(events.APIGatewayProxyRequest{
		RequestContext: requestCtx,
		Body:           string(requestBytes),
	}); err != nil {
		t.Error(err)
	} else if r.StatusCode != 200 {
		t.Error(fmt.Errorf("bad response [%v]", r))
	}
}

// tests an invalid save, no body or token
func TestHandleRequestSaveBadToken401(t *testing.T) {
	addressRequest.AccessToken = ""
	addressRequest.Command = "save"
	requestBytes, _ := json.Marshal(&addressRequest)
	if r, err := HandleRequest(events.APIGatewayProxyRequest{
		Body: string(requestBytes),
	}); err != nil {
		t.Error(err)
	} else if r.StatusCode != 401 {
		t.Error(fmt.Errorf("bad response [%v]", r))
	}
}

// tests the golden path for the find by ids command
func TestHandleRequestFindByIds200(t *testing.T) {
	addressRequest.Command = "find-by-ids"
	requestBytes, _ := json.Marshal(&addressRequest)
	if r, err := HandleRequest(events.APIGatewayProxyRequest{
		RequestContext:        requestCtx,
		Body:                  string(requestBytes),
		QueryStringParameters: map[string]string{"cmd": "find-by-ids", "ids": addressId},
	}); err != nil {
		t.Error(err)
	} else if r.StatusCode != 200 {
		t.Error(fmt.Errorf("bad response [%v]", r))
	}
}

// tests the golden path for the find by ids command
func TestHandleRequestFindByIds400(t *testing.T) {
	addressRequest.Command = "find-by-ids"
	addressRequest.Ids = nil
	requestBytes, _ := json.Marshal(&addressRequest)
	if r, err := HandleRequest(events.APIGatewayProxyRequest{
		RequestContext: requestCtx,
		Body:           string(requestBytes),
	}); err != nil {
		t.Error(err)
	} else if r.StatusCode != 400 {
		t.Error(fmt.Errorf("bad response [%v]", r))
	}
}

// tests default switch case, a bad request
func TestHandleRequestBadCommand400(t *testing.T) {
	addressRequest.Command = ""
	if r, err := HandleRequest(events.APIGatewayProxyRequest{}); err != nil {
		t.Error(err)
	} else if r.StatusCode != 400 {
		t.Error(fmt.Errorf("bad response [%v]", r))
	}
}

// for code coverage purposes only
func TestHandleRequest(t *testing.T) {
	go main()
}
