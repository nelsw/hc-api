package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/internal/entity/token"
)

func Handle(t token.Aggregate) (string, error) {
	if err := t.Validate(); err != nil {
		return "", err
	} else {
		return t.String(), nil
	}
}

func main() {
	lambda.Start(Handle)
}
