package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	. "hc-api/service"
	"os"
	"strings"
)

var t = os.Getenv("ADDRESS_TABLE")

type Address struct {
	Id     string `json:"id"`
	Street string `json:"street"`
	Unit   string `json:"unit,omitempty"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip5   string `json:"zip_5"`
	Zip4   string `json:"zip_4,omitempty"`
}

type AddressRequest struct {
	AccessToken string   `json:"access_token"`
	Command     string   `json:"command"`
	Ids         []string `json:"ids"`
	Address     Address  `json:"address"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Printf("REQUEST   [%s]\n", request.Body)

	body := request.Body
	cmd := request.QueryStringParameters["cmd"]
	session := request.QueryStringParameters["session"]
	ids := strings.Split(request.QueryStringParameters["ids"], ",")

	var v AddressRequest
	if err := json.Unmarshal([]byte(body), &v); err == nil {
		b, _ := json.Marshal(v.Address)
		cmd = v.Command
		session = v.AccessToken
		ids = v.Ids
		body = string(b)
	}

	switch cmd {

	case "save":
		var a Address
		if _, err := ValidateSession(session, request.RequestContext.Identity.SourceIP); err != nil {
			return Unauthorized().Error(err).Build()
		} else {
			str, _ := VerifyAddress(body)
			_ = json.Unmarshal([]byte(str), &a)
			_ = Put(a, &t)
			return Ok().Data(&a).Build()
		}

	case "find-by-ids":
		var p []Address
		if err := FindAllById(t, ids, &p); err != nil {
			return BadRequest().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}

	default:
		return BadRequest().Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
