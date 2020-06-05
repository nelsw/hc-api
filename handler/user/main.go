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
	} else if token, ok := r.Headers["token"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no token"))
	} else if id, ok := r.QueryStringParameters["id"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no id"))
	} else if csv, ok := r.QueryStringParameters["ids"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no ids"))
	} else if col, ok := r.QueryStringParameters["col"]; !ok {
		return apigwp.Response(400, fmt.Errorf("no col"))
	} else {

		authenticate := events.APIGatewayProxyRequest{Path: "authenticate", Headers: r.Headers}
		if code, body := client.CallIt(authenticate, "tokenHandler"); code != 200 {
			return apigwp.Response(code, body)
		} else {
			r.Headers["Authorize"] = body
		}

		keyword := r.Path + " " + col
		ids := strings.Split(csv, ",")
		m := map[string]interface{}{"table": table, "id": id, "ids": ids, "keyword": keyword}
		return apigwp.Response(client.CallIt(&m, "repoHandler"))
	}
}

func main() {
	lambda.Start(Handle)
}
