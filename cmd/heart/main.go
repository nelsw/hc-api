package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"sam-app/pkg/factory/apigwp"
)

var m = map[string]string{
	"address_api": os.Getenv("ADDRESS_API_GW_URL"),
	"offer_api":   os.Getenv("OFFER_API_GW_URL"),
	"order_api":   os.Getenv("ORDER_API_GW_URL"),
	"product_api": os.Getenv("PRODUCT_API_GW_URL"),
	"profile_api": os.Getenv("PROFILE_API_GW_URL"),
	"user_api":    os.Getenv("USER_API_GW_URL"),
}

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQUEST   [%v]", request)
	return apigwp.Response(200, m)
}

func main() {
	lambda.Start(Handle)
}
