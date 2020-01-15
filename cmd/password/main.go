package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/internal/entity/password"
	. "hc-api/service"
)

func Handle(p password.Password) error {
	if err := FindOne(p.Table(), p.Id(), &p); err != nil {
		return err
	} else {
		return p.Validate()
	}
}

func main() {
	lambda.Start(Handle)
}
