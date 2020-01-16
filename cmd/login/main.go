package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
	"hc-api/pkg/factory"
	"hc-api/pkg/service"
)

// When provided a valid proxy wrapper body, this Handler returns a 24 hour JWT access token.
func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Credential composes client and server c to satisfy the Password interface.
	c := entity.Credential{}
	if err := factory.Request(r, &c); err != nil {
		// Unable to unmarshal or validate ClientCredential.
		return factory.Response(400, err)
	} else if err := service.Find(&c); err != nil {
		// ServerCredential referred to by ClientCredential not found.
		return factory.Response(404, err)
	}

	// Credentials are valid, let's find the p associated with the server c.
	// Password entity found, but does it match the encoded server digest value?
	p := entity.Password{Id: c.PasswordId}
	if _, err := service.Invoke(&p); err != nil {
		// Passwords did not match.
		return factory.Response(401, err)
	}

	// Passwords matched, lets create and return an access token.
	a := entity.Authorization{UserId: c.UserId, SourceIp: r.RequestContext.Identity.SourceIP}
	t := entity.Token{Authorization: a}
	if str, err := service.Invoke(&t); err != nil {
		// Password could not be processed, theoretically impossible when creating tokens.
		return factory.Response(422, err)
	} else {
		// Password successfully aggregated, authorized, and validated; 24 hour access token returned.
		return factory.Response(200, str)
	}
}

func main() {
	lambda.Start(Handle)
}
