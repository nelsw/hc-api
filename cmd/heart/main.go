package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
)

var headers = map[string]string{"Access-Control-Allow-Origin": "*"}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQUEST [%v]", request)
	b, _ := json.Marshal(map[string]string{
		"address_api": os.Getenv("ADDRESS_API_GW_URL"),
		"offer_api":   os.Getenv("OFFER_API_GW_URL"),
		"order_api":   os.Getenv("ORDER_API_GW_URL"),
		"product_api": os.Getenv("PRODUCT_API_GW_URL"),
		"profile_api": os.Getenv("PROFILE_API_GW_URL"),
		"user_api":    os.Getenv("USER_API_GW_URL"),
	})
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(b),
	}
	fmt.Printf("RESPONSE [%v]", response)
	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
