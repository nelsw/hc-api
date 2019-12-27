package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	. "hc-api/service"
	"os"
)

type Offer struct {
	Id           string `json:"id"`
	UserId       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	ProductId    string `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductPrice int64  `json:"product_price"`
	ProductQty   int    `json:"product_qty"`
	ProductImg   string `json:"product_img,omitempty"`
	OfferTotal   int64  `json:"offer_total"`
	Created      string `json:"created"`
}

var tableName = os.Getenv("OFFER_TABLE")

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("REQUEST [%v]", request)

	var statusCode = 400
	var body = "bad request"

	if request.QueryStringParameters["cmd"] == "save" {
		var offer Offer
		if err := json.Unmarshal([]byte(request.Body), &offer); err != nil {
			statusCode = 502
			body = err.Error()
		} else if err := Put(offer, &tableName); err != nil {
			statusCode = 500
			body = err.Error()
		} else {
			statusCode = 200
			b, _ := json.Marshal([]byte(`{"message":"success"}`))
			body = string(b)
		}
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       body,
	}

	fmt.Printf("RESPONSE [%v]", response)

	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
