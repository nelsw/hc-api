package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
	"hc-api/pkg/service"
)

func HandleRequest(user entity.User) error {
	if user.Id != "" {
		return service.Find(&user)
	} else {
		return service.Update(&user)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
