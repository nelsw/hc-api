package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/cmd/shipping/fedex"
	"hc-api/cmd/shipping/ups"
	"hc-api/cmd/shipping/usps"
	"hc-api/pkg/factory"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	ip := r.RequestContext.Identity.SourceIP
	body := r.Body
	fmt.Printf("REQUEST [%s]: ip=[%s], body=[%s]\n", cmd, ip, body)
	switch r.QueryStringParameters["cmd"] {
	case "verify":
		if a, err := usps.GetAddress(body); err != nil {
			return factory.Response(400, err)
		} else {
			return factory.Response(200, &a)
		}
	case "rate":
		v := r.QueryStringParameters["v"]
		if v == "USPS" {
			if p, err := usps.GetPostage(body); err != nil {
				return factory.Response(400, err)
			} else {
				return factory.Response(200, &p)
			}
		} else if v == "UPS" {
			if p, err := ups.GetRates(body); err != nil {
				return factory.Response(400, err)
			} else {
				return factory.Response(200, &p)
			}
		} else if v == "FEDEX" {
			if p, err := fedex.GetRates(body); err != nil {
				return factory.Response(400, err)
			} else {
				return factory.Response(200, &p)
			}
		}
	}
	return factory.Response(400)
}

func main() {
	lambda.Start(HandleRequest)
}
