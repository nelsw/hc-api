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
	log.Printf("request=%v", r)

	cmd := r.QueryStringParameters["cmd"]

	switch cmd {

	case "save":
		var a model.Address
		if err := json.Unmarshal([]byte(r.Body), &a); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := a.Validate(); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if err := a.SetId(); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := repo.SaveAddress(&a); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&a).Build()
		}

	case "find":
		var a model.Address
		if err := json.Unmarshal([]byte(r.Body), &a); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := repo.FindAddress(&a); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&a).Build()
		}

	case "findBatch":
		var aa []model.Address
		if err := json.Unmarshal([]byte(r.Body), &aa); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else if err := repo.FindAllAddresses(&aa); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&aa).Build()
		}

	case "create":
		return response.New().Code(http.StatusNotImplemented).Build()

	case "update":
		return response.New().Code(http.StatusNotImplemented).Build()

	case "delete":
		return response.New().Code(http.StatusNotImplemented).Build()

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
