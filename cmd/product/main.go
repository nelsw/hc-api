package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	response "github.com/nelsw/hc-util/aws"
	"hc-api/model"
	"hc-api/repo"
	"log"
	"net/http"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("request: [%v]", r)

	cmd := r.QueryStringParameters["cmd"]

	switch cmd {

	case "save":
		var p model.Product
		if err := json.Unmarshal([]byte(r.Body), &p); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := p.Validate(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := repo.SaveProduct(&p); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&p).Build()
		}

	case "find-all":
		if products, err := repo.FindAllProducts(); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&products).Build()
		}

	case "find-by-owner":
		s := r.QueryStringParameters["owner"]
		if products, err := repo.FindAllProductsByOwner(&s); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&products).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
