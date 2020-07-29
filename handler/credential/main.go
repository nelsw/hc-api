package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/factory/apigwp"
	"sam-app/pkg/model/credential"
)

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	e := credential.Entity{}

	ip, err := apigwp.Request(r, &e)
	if err != nil {
		return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: 400, Headers: r.Headers, Body: err.Error()})
	}

	p := []byte(e.Password)
	i := map[string]interface{}{"Table": "credential", "Type": "*credential.Entity", "Keyword": "find-one", "Id": e.Id}
	code, body := client.CallIt(i, "repoHandler")
	_ = json.Unmarshal([]byte(body), &e)

	if e.UserId == "" {
		code = 404
	}

	if code != 200 {
		return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: code, Headers: r.Headers, Body: body})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(e.Password), p); err != nil {
		return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: 401, Headers: r.Headers, Body: err.Error()})
	}

	claims := jwt.StandardClaims{ip, 0, e.UserId, 0, "credentialHandler", 0, "login"}
	b, _ := json.Marshal(&claims)

	authRequest := events.APIGatewayProxyRequest{Path: "authorize", Body: string(b)}
	return apigwp.HandleResponse(client.Invoke("tokenHandler", authRequest))
}

func main() {
	lambda.Start(Handle)
}
