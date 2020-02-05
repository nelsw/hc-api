package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/client/faas/client"
	"hc-api/pkg/client/repo"
	"hc-api/pkg/factory/apigwp"
	"hc-api/pkg/model/credential"
	"hc-api/pkg/model/password"
	"hc-api/pkg/model/token"
	"strings"
)

// This handler returns returns a 24 hour JWT access token when provided with a valid Value validation.
func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var c credential.Entity

	ip, err := apigwp.Request(r, &c)
	if err != nil {
		return apigwp.Response(400, err)
	}

	p := password.Entity{c.PasswordId, ""}
	if err := repo.FindOne(&p); err != nil || p.Hash == "" {
		return apigwp.Response(404, err)
	} else if err := p.ComparePasswords(c.Password); err != nil {
		return apigwp.Response(401, err)
	}

	t := token.Value{c.UserId, ip, nil, "authorize"}
	out, err := client.Invoke(&token.Entity{t, token.Error{}})
	if err != nil || strings.Contains(string(out), "errorMessage") {
		return apigwp.Response(500, err)
	} else {
		return apigwp.Response(200, out)
	}
}

func main() {
	lambda.Start(Handle)
}
