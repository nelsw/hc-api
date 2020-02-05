package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/model/fedex"
)

func Handle(i fedex.RateRequest) (interface{}, error) {
	return fedex.GetRates(i)
}

func main() {
	lambda.Start(Handle)
}
