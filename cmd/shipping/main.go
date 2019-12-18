package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/cmd/shipping/usps"
	. "hc-api/service"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQUEST [%v]\n", r)
	switch r.QueryStringParameters["cmd"] {
	case "verify":
		if a, err := usps.GetAddress(r.Body); err != nil {
			return BadRequest().Error(err).Build()
		} else {
			return Ok().Data(&a).Build()
		}
	case "rate":
		fmt.Printf(r.Body)
		if p, err := usps.GetPostage(r.Body); err != nil {
			return BadRequest().Error(err).Build()
		} else {
			return Ok().Data(&p).Build()
		}
	default:
		return BadRequest().Data(r).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
