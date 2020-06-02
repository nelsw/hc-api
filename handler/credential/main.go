package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/client/repo"
	"sam-app/pkg/factory/apigwp"
	"sam-app/pkg/model/credential"
)

var table = os.Getenv("TABLE")

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var e credential.Entity

	ip, err := apigwp.Request(r, &e)
	if err != nil {
		return apigwp.Response(400, err)
	}

	p := e.Password
	if err := repo.FindById(table, e.Id, &e); err != nil {
		return apigwp.Response(500, err)
	} else if err := bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(p)); err != nil {
		return apigwp.Response(401, err)
	}

	claims := jwt.StandardClaims{
		ip,
		0,
		e.Id,
		0,
		"credentialHandler",
		0,
		"login",
	}
	b, _ := json.Marshal(&claims)

	authorizationRequest := events.APIGatewayProxyRequest{Path: "authorize", Body: string(b)}

	return apigwp.Response(client.CallIt(authorizationRequest, "tokenHandler"))
}

func main() {
	lambda.Start(Handle)
}
