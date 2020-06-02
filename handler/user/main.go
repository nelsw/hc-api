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
	"strings"
)

var table = os.Getenv("TABLE")

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apigwp.LogRequest(r)
	if r.Path != "add" && r.Path != "delete" {
		return apigwp.Response(400, fmt.Errorf("bad path [%s]", r.Path))
	} else if token, ok := r.QueryStringParameters["token"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no token provided"))
	} else if csv, ok := r.QueryStringParameters["ids"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no ids provided"))
	} else if col, ok := r.QueryStringParameters["col"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no col provided"))
	} else {

		keyword := r.Path + " " + col
		inspection := events.APIGatewayProxyRequest{
			Path:                  "inspect",
			QueryStringParameters: map[string]string{"token": token},
		}
		code, body := client.CallIt(inspection, "tokenHandler")
		if code != 200 {
			return apigwp.Response(code, body)
		}

		response := events.APIGatewayProxyResponse{}
		_ = json.Unmarshal([]byte(body), &response)

		if response.StatusCode != 200 {
			return apigwp.Response(response.StatusCode, response.Body)
		}

		claims := jwt.StandardClaims{}
		if err := json.Unmarshal([]byte(response.Body), &claims); err != nil {
			return apigwp.Response(400, err)
		}

		ids := strings.Split(csv, ",")
		m := map[string]interface{}{
			"table":   table,
			"id":      claims.Id,
			"ids":     ids,
			"keyword": keyword,
		}
		return apigwp.Response(client.CallIt(&m, "repoHandler"))
	}
}

func main() {
	lambda.Start(Handle)
}
