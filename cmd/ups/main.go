package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/model/ups"
)

func Handle(r ups.PostageRateRequest) (interface{}, error) {
	return ups.GetRates(r)
}

func main() {
	lambda.Start(Handle)
}
