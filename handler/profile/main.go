package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/factory/apigwp"
	"sam-app/pkg/model/profile"
)

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var e profile.Entity

	if _, err := apigwp.Request(r, &e); err != nil {
		return apigwp.Response(400, err)
	}

	authResponse := client.Invoke("tokenHandler", events.APIGatewayProxyRequest{Path: "inspect", Headers: r.Headers})
	if authResponse.StatusCode != 200 {
		return apigwp.HandleResponse(authResponse)
	}

	claims := map[string]interface{}{"jti": ""}
	_ = json.Unmarshal([]byte(authResponse.Body), &claims)

	switch r.QueryStringParameters["path"] {

	case "save":

		m := map[string]interface{}{"id": claims["jti"], "table": "profile", "type": "*profile.Entity", "keyword": "save", "result": &e}
		return apigwp.Response(client.CallIt(&m, "repoHandler"))
	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(Handle)
}
