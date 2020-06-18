package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	"os"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/factory/apigwp"
	"sam-app/pkg/model/profile"
)

var table = os.Getenv("PROFILE_TABLE")

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var e profile.Entity

	if _, err := apigwp.Request(r, &e); err != nil {
		return apigwp.Response(400, err)
	}

	authenticate := events.APIGatewayProxyRequest{Path: "authenticate", Headers: r.Headers}
	authResponse := client.Invoke("tokenHandler", authenticate)
	if authResponse.StatusCode != 200 {
		return apigwp.Response(authResponse.StatusCode, authResponse.Body)
	}
	claims := jwt.StandardClaims{}
	_ = json.Unmarshal([]byte(authResponse.Body), &claims)

	switch r.Path {

	case "save":

		m := map[string]interface{}{"id": claims.Id, "table": table, "type": "*profile.Entity", "keyword": "save", "result": &e}
		return apigwp.Response(client.CallIt(&m, "repoHandler"))
	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(Handle)
}
