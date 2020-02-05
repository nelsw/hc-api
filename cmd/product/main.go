package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"hc-api/pkg/client/faas/client"
	"hc-api/pkg/client/repo"
	"hc-api/pkg/factory/apigwp"
	"hc-api/pkg/model/product"
	"hc-api/pkg/model/token"
)

var InvalidRequest = fmt.Sprintf("bad request\n")
var UnauthorizedRequest = fmt.Sprintf("token invalid or expired\n")

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var e product.Proxy

	ip, err := apigwp.Request(r, &e)
	if err != nil {
		return apigwp.Response(400, err)
	}

	op := e.Subject

	e.Subject = "authenticate"
	e.SourceIp = ip

	tkn := token.Entity{e.Value, token.Error{}}

	if out, err := client.Invoke(&tkn); err != nil {
		return apigwp.Response(500, err)
	} else if err := json.Unmarshal(out, &tkn); err != nil {
		return apigwp.Response(500, err)
	}

	if tkn.SourceId == "" {
		return apigwp.Response(402, UnauthorizedRequest)
	}

	e.OwnerId = tkn.SourceId

	switch op {

	case "save":
		newProduct := e.Id == ""
		if newProduct {
			s, _ := uuid.NewUUID()
			e.Id = s.String()
		}
		if err := repo.SaveOne(&e.Entity); err != nil {
			return apigwp.Response(500, err)
		}
		if newProduct {
			// todo - update user
		}
		return apigwp.Response(200)

	case "find-one":
		if err := repo.FindOne(&e.Entity); err != nil {
			return apigwp.Response(404, err)
		}
		return apigwp.Response(200, &e)

	case "find-many":
		if out, err := repo.FindMany(&e.Entity, e.Ids); err != nil {
			return apigwp.Response(404, err)
		} else {
			return apigwp.Response(200, &out)
		}

	}

	return apigwp.Response(400, InvalidRequest)
}

func main() {
	lambda.Start(Handle)
}
