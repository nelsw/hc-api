package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/cmd/shipping/fedex"
	"hc-api/cmd/shipping/ups"
	"hc-api/cmd/shipping/usps"
	. "hc-api/service"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	ip := r.RequestContext.Identity.SourceIP
	body := r.Body
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]\n", cmd, ip, body)
	switch r.QueryStringParameters["cmd"] {
	case "verify":
		if a, err := usps.GetAddress(body); err != nil {
			return BadRequest().Error(err).Build()
		} else {
			return Ok().Data(&a).Build()
		}
	case "rate":
		fmt.Printf(r.Body)
		v := r.QueryStringParameters["v"]
		if v == "USPS" {
			if p, err := usps.GetPostage(body); err != nil {
				return BadRequest().Error(err).Build()
			} else {
				return Ok().Data(&p).Build()
			}
		} else if v == "UPS" {
			if p, err := ups.GetRates(body); err != nil {
				return BadRequest().Error(err).Build()
			} else {
				return Ok().Data(&p).Build()
			}
		} else if v == "FEDEX" {
			if p, err := fedex.GetRates(body); err != nil {
				return BadRequest().Error(err).Build()
			} else {
				return Ok().Data(&p).Build()
			}
		}
	}
	return BadRequest().Data(r).Build()
}

func main() {
	lambda.Start(HandleRequest)
}
