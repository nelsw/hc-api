// UserProfile is exactly what it appears to be, a user profile domain model entity.
// It also promotes separation of concerns by decoupling user profile details from the primary user entity. IF
// UserProfile.EmailOld != UserProfile.EmailNew, AND User.Email == UserProfile.EmailOld, THEN we must prompt the user to
// confirm new email address. IF UserProfile.Password1 is not blank AND UserProfile.Password2 is not blank AND valid AND
// UserProfile.Password1 == UserProfile.Password2, then we update the UserPassword entity and return OK.
package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
	"hc-api/pkg/factory"
	"hc-api/pkg/service"
)

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	e := entity.Profile{}
	if err := factory.Request(r, &e); err != nil {
		return factory.Response(400, err)
	}

	e.Authorization.SourceIp = r.RequestContext.Identity.SourceIP
	t := entity.Token{Authorization: e.Authorization}
	if _, err := service.Invoke(&t); err != nil {
		return factory.Response(402, err)
	}

	switch e.Case {

	case "save":
		if err := service.Save(&e); err != nil {
			return factory.Response(400, err)
		} else {
			return factory.Response(200, &e)
		}

	case "find":
		if err := service.Find(&e); err != nil {
			return factory.Response(400, err)
		} else {
			return factory.Response(200, &e)
		}
	}

	return factory.Response(400, fmt.Sprintf("bad case=[%s]", e.Case))
}

func main() {
	lambda.Start(Handle)
}
