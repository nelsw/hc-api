package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	response "github.com/nelsw/hc-util/aws"
	"hc-api/model"
	"hc-api/service"
	"log"
	"net/http"
	"strings"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("request: [%v]", r)

	cmd := r.QueryStringParameters["cmd"]

	switch cmd {

	case "save":
		var a model.Address
		if err := a.Unmarshal(r.Body); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if _, err := service.ValidateSession(a.Session, r.RequestContext.Identity.SourceIP); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := service.SaveAddress(&a); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&a).Build()
		}

	case "find-by-ids":
		csv := r.QueryStringParameters["ids"]
		ids := strings.Split(csv, ",")
		if aa, err := service.FindAllAddressesByIds(&ids); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&aa).Build()
		}

	default:
		return response.New().Code(http.StatusBadRequest).Text(fmt.Sprintf("bad command: [%s]", cmd)).Build()
	}
}

func main() {
	lambda.Start(HandleRequest)
}
