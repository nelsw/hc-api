package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/internal/entity/offer"
	"hc-api/internal/factory/apigwp"
	"hc-api/internal/service"
	. "hc-api/service"
)

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var o offer.Aggregate
	if err := apigwp.Request(r, &o); err != nil {
		// Unable to unmarshal or validate ClientCredential.
		return apigwp.Response(400, err)
	}

	userId, err := service.String(&o.Token)
	if err != nil {
		// Entity could not be processed, theoretically impossible when creating tokens.
		return apigwp.Response(422, err)
	}

	if err := service.Find(&o.User, userId); err != nil {
		return apigwp.Response(404, err)
	}

	if err := service.Find(&o.Profile, o.ProfileId); err != nil {
		return apigwp.Response(404, err)
	}

	o.Object.SetNewId()
	o.UserId = userId

	if err := Put(o, o.Name()); err != nil {
		return apigwp.Response(400, err)
	} else {
		return apigwp.Response(200)
	}
}

func main() {
	lambda.Start(Handle)
}
