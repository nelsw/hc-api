package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
	"hc-api/pkg/factory"
	"hc-api/pkg/service"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	e := entity.Product{}
	if err := factory.Request(r, &e); err != nil {
		return factory.Response(400, err)
	}

	e.SourceIp = r.RequestContext.Identity.SourceIP
	t := entity.Token{Authorization: e.Authorization}
	uid, err := service.Invoke(&t)
	if err != nil {
		return factory.Response(422, err)
	}

	e.UserId = string(uid)
	e.OwnerId = e.UserId
	e.Stock = e.Quantity
	if out, err := service.Invoke(&e); err != nil {
		return factory.Response(400, err)
	} else {
		return factory.Response(200, &out)
	}

}

func main() {
	lambda.Start(HandleRequest)
}
