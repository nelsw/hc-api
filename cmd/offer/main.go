package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/client/faas/client"
	"hc-api/pkg/client/repo"
	"hc-api/pkg/factory/apigwp"
	"hc-api/pkg/model/offer"
	"hc-api/pkg/model/token"
)

var InvalidRequest = fmt.Sprintf("bad request\n")

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var e offer.Proxy

	ip, err := apigwp.Request(r, &e)
	if err != nil {
		return apigwp.Response(400, err)
	}

	op := e.Subject

	e.Subject = "authenticate"
	e.SourceIp = ip

	tkn := token.Entity{e.Value, token.Error{}}
	out, err := client.Invoke(&tkn)
	if err != nil {
		return apigwp.Response(500, err)
	}

	_ = json.Unmarshal(out, &tkn)
	if tkn.Error.Msg != "" {
		return apigwp.Response(402, &tkn)
	}

	e.UserId = tkn.SourceId

	switch op {

	case "save":
		if err := repo.SaveOne(&e.Entity); err != nil {
			return apigwp.Response(500, err)
		}
		return apigwp.Response(200)
	}

	return apigwp.Response(400, InvalidRequest)

}

func main() {
	lambda.Start(Handle)
}
