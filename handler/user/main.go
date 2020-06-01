package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"sam-app/pkg/factory/apigwp"
)

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return apigwp.Response(400, "no op")
}

func main() {
	lambda.Start(Handle)
}
