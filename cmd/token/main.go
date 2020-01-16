package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
)

func Handle(e entity.Token) ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	} else {
		return e.Payload(), nil
	}
}

func main() {
	lambda.Start(Handle)
}
