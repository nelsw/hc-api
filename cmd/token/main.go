package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/model/token"
)

func Handle(e token.Entity) (interface{}, error) {
	fmt.Printf("REQUEST  entity=[%v]\n", e)
	if err := e.Validate(); err != nil {
		return nil, err
	}

	if e.Subject == "authenticate" {
		err := e.Authenticate()
		return e, err
	}

	if e.Subject == "authorize" {
		err := e.Authorize()
		return e, err
	}

	return nil, token.InvalidToken
}

func main() {
	lambda.Start(Handle)
}
