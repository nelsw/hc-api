package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
	"hc-api/pkg/service"
)

func Handle(p entity.Password) error {
	if err := service.Find(&p); err != nil {
		return err
	} else {
		return p.Validate()
	}
}

func main() {
	lambda.Start(Handle)
}
