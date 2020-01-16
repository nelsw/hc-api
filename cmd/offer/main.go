package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
	"hc-api/pkg/factory"
	"hc-api/pkg/service"
)

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	e := entity.Offer{}
	if err := factory.Request(r, &e); err != nil {
		// Unable to unmarshal or validate ClientCredential.
		return factory.Response(400, err)
	}

	e.Authorization.SourceIp = r.RequestContext.Identity.SourceIP
	t := entity.Token{Authorization: e.Authorization}
	if uid, err := service.Invoke(&t); err != nil {
		// Entity could not be processed, theoretically impossible when creating tokens.
		return factory.Response(422, err)
	} else {
		e.UserId = string(uid)
	}

	u := entity.User{Id: e.UserId}
	if err := service.Find(&u); err != nil {
		return factory.Response(404, err)
	} else {
		e.ProfileId = u.ProfileId
	}

	p := entity.Profile{Id: e.ProfileId}
	if err := service.Find(&p); err != nil {
		return factory.Response(404, err)
	} else {
		e.Phone = p.Phone
		e.Email = p.Email
		e.FirstName = p.FirstName
		e.LastName = p.LastName
	}

	if out, err := service.Invoke(&e); err != nil {
		return factory.Response(400, err)
	} else {
		return factory.Response(200, &out)
	}
}

func main() {
	lambda.Start(Handle)
}
