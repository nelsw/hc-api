package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
)

func Handle(t entity.Token) error {
	return t.Validate()
}

func main() {
	lambda.Start(Handle)
}
