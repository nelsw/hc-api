package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	response "github.com/nelsw/hc-util/aws"
	"hc-api/model"
	"hc-api/service"
	"net/http"
	"strings"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cmd := r.QueryStringParameters["cmd"]
	fmt.Printf("REQUEST [%s]: [%v]", cmd, r)

	switch cmd {

	case "save":
		var p model.Product
		if err := p.Unmarshal(r.Body); err != nil {
			return response.New().Code(http.StatusBadRequest).Text(err.Error()).Build()
		} else if _, err := service.ValidateSession(p.Session, r.RequestContext.Identity.SourceIP); err != nil {
			return response.New().Code(http.StatusUnauthorized).Text(err.Error()).Build()
		} else if err := service.SaveProduct(&p); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&p).Build()
		}

	case "find-all":
		if products, err := service.FindAllProducts(); err != nil {
			return response.New().Code(http.StatusInternalServerError).Text(err.Error()).Build()
		} else {
			return response.New().Code(http.StatusOK).Data(&products).Build()
		}

	case "find-by-ids":
		csv := r.QueryStringParameters["ids"]
		ids := strings.Split(csv, ",")
		if products, err := service.FindAllProductsByIds(&ids); err != nil {
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
