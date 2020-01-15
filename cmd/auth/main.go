package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/internal/entity/token"
)

func Handle(t token.Aggregate) error {
	return t.Validate()
}

func main() {
	lambda.Start(Handle)
}
