package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	response "github.com/nelsw/hc-util/aws"
	"hc-api/repo"
	"log"
	"net/http"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("request: [%v]", r)

	cmd := r.QueryStringParameters["cmd"]

	switch cmd {

	case "find-all":
		if products, err := repo.FindAllProducts(); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(products).Build()
		}

	case "create":
		return response.New().Code(http.StatusNotImplemented).Build()

	case "update":
		return response.New().Code(http.StatusNotImplemented).Build()

	case "delete":
		return response.New().Code(http.StatusNotImplemented).Build()

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
