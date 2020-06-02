package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/factory/apigwp"
)

var table = os.Getenv("TABLE")

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if r.Path == "add" || r.Path == "delete" {
		m := map[string]interface{}{
			"table":   table,
			"ids":     r.QueryStringParameters["ids"],
			"keyword": r.Path + " " + r.QueryStringParameters["col"],
		}
		return apigwp.Response(client.CallIt(&m, "repoHandler"))
	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(Handle)
}
