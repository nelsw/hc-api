package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/internal/entity/credential"
	"hc-api/internal/entity/password"
	"hc-api/internal/entity/token"
	"hc-api/internal/factory/apigwp"
	"hc-api/internal/service"
)

// When provided a valid proxy wrapper body, this Handler returns a 24 hour JWT access token.
func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Aggregate composes client and server credential to satisfy the Entity interface.
	var a credential.Aggregate
	if err := apigwp.Request(r, &a); err != nil {
		// Unable to unmarshal or validate ClientCredential.
		return apigwp.Response(400, err)
	} else if err := service.Find(&a); err != nil {
		// ServerCredential referred to by ClientCredential not found.
		return apigwp.Response(404, err)
	}

	// Credentials are valid, let's find the password associated with the server credential.
	p := &password.Entity{ID: a.PasswordId, Decoded: a.Payload()}
	if err := service.Find(p); err != nil {
		// Password entity referred to by ClientCredential entity not found.
		return apigwp.Response(404, err)
	}

	// Password entity found, but does it match the encoded server digest value?
	if err := service.Validate(p); err != nil {
		// Passwords did not match.
		return apigwp.Response(401, err)
	}

	// Passwords matched, lets create and return an access token.
	t := &token.Aggregate{UserId: a.UserId, SourceIp: r.RequestContext.Identity.SourceIP}
	if str, err := service.String(t); err != nil {
		// Entity could not be processed, theoretically impossible when creating tokens.
		return apigwp.Response(422, err)
	} else {
		// Entity successfully aggregated, authorized, and validated; 24 hour access token returned.
		return apigwp.Response(200, str)
	}
}

func main() {
	lambda.Start(Handle)
}
