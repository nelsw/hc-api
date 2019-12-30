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
	Id      string `json:"id"`
	Session string `json:"session"`
	Street  string `json:"street"`
	Unit    string `json:"unit,omitempty"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip5    string `json:"zip_5"`
	Zip4    string `json:"zip_4,omitempty"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQUEST [%v]", request)

	switch request.QueryStringParameters["cmd"] {

	case "save":
		var a Address
		session := request.QueryStringParameters["session"]
		if _, err := ValidateSession(session, request.RequestContext.Identity.SourceIP); err != nil {
			return Unauthorized().Error(err).Build()
		} else if str, err := VerifyAddress(request.Body); err != nil {
			return InternalServerError().Error(err).Build()
		} else if err := json.Unmarshal([]byte(str), &a); err != nil {
			return InternalServerError().Error(err).Build()
		} else if err := Put(a, &t); err != nil {
			return InternalServerError().Error(err).Build()
		} else {
			return Ok().Data(&a).Build()
		}

	case "find-by-ids":
		var p []Address
		ss := strings.Split(request.QueryStringParameters["ids"], ",")
		if err := FindAllById(t, ss, &p); err != nil {
			return InternalServerError().Error(err).Build()
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
