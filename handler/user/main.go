package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	} else if _, ok := r.Headers["Authorization"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no token"))
	} else if id, ok := r.QueryStringParameters["id"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no id"))
	} else if csv, ok := r.QueryStringParameters["ids"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no ids"))
	} else if col, ok := r.QueryStringParameters["col"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no col"))
	} else {

		authenticate := events.APIGatewayProxyRequest{Path: "authenticate", Headers: r.Headers}
		if authResponse := client.Invoke("tokenHandler", authenticate); authResponse.StatusCode != 200 {
			return apigwp.Response(401, authResponse.Body)
		} else {
			r.Headers = authResponse.Headers
		}

		keyword := r.Path + " " + col
		ids := strings.Split(csv, ",")
		m := map[string]interface{}{"table": table, "id": id, "ids": ids, "keyword": keyword, "type": "*user.Entity"}
		code, body := client.CallIt(&m, "repoHandler")
		return apigwp.ProxyResponse(code, r.Headers, body)
	}
}

func main() {
	lambda.Start(Handle)
}
