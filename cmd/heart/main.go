package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/internal/factory/apigwp"
	"os"
	"strings"
)

func ApiUrls(s string) map[string]string {
	if s != "" {
		s = "_" + strings.ToUpper(s)
	}
	return map[string]string{
		"address_api": os.Getenv("ADDRESS_API_GW_URL" + s),
		"offer_api":   os.Getenv("OFFER_API_GW_URL" + s),
		"order_api":   os.Getenv("ORDER_API_GW_URL" + s),
		"product_api": os.Getenv("PRODUCT_API_GW_URL" + s),
		"profile_api": os.Getenv("PROFILE_API_GW_URL" + s),
		"user_api":    os.Getenv("USER_API_GW_URL" + s),
	}
}

func Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQUEST   [%v]", request)
	return apigwp.Response(200, ApiUrls(request.QueryStringParameters["sd"]))
}

func main() {
	lambda.Start(Handle)
}
