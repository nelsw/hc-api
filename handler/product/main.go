package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"sam-app/pkg/client/repo"
	"sam-app/pkg/factory/apigwp"
	"sam-app/pkg/model/product"
	"strings"
)

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	apigwp.LogRequest(r)

	switch r.Path {

	case "save":
		p := product.Entity{}
		if err := json.Unmarshal([]byte(r.Body), &p); err != nil {
			return apigwp.Response(400, err)
		} else if err := repo.SaveOne(&p); err != nil {
			return apigwp.Response(500, err)
		} else {
			return apigwp.Response(200, &p)
		}

	case "remove":
		if id, ok := r.QueryStringParameters["id"]; ok {
			p := product.Entity{}
			if err := repo.Remove(&p, id); err != nil {
				return apigwp.Response(500, err)
			} else {
				return apigwp.Response(200, &p)
			}
		}

	case "find":
		if csv, ok := r.QueryStringParameters["ids"]; ok {
			ids := strings.Split(csv, ",")
			if out, err := repo.FindMany(&product.Entity{}, ids); err != nil {
				return apigwp.Response(500, err)
			} else {
				return apigwp.Response(200, &out)
			}
		} else if id, ok := r.QueryStringParameters["id"]; ok {
			p := product.Entity{Id: id}
			if err := repo.FindOne(&p); err != nil {
				return apigwp.Response(500, err)
			} else {
				return apigwp.Response(200, &p)
			}
		}

	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(Handle)
}
