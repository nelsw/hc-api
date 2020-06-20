package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/factory/apigwp"
	"sam-app/pkg/model/product"
	"strings"
)

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	apigwp.LogRequest(r)

	switch r.QueryStringParameters["path"] {

	case "save":
		p := product.Entity{}

		if _, err := apigwp.Request(r, &p); err != nil {
			return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: 400, Headers: r.Headers, Body: err.Error()})
		}

		i := map[string]interface{}{"Table": "product", "Type": "*product.Entity", "Keyword": "save", "Id": p.Id, "Result": &p}
		if code, body := client.CallIt(i, "repoHandler"); code != 200 {
			return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: code, Headers: r.Headers, Body: body})
		} else if b, err := json.Marshal(&p); err != nil {
			return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: 500, Headers: r.Headers, Body: body})
		} else {
			return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: code, Headers: r.Headers, Body: string(b)})
		}

	case "remove":
		if id, ok := r.QueryStringParameters["id"]; ok {
			i := map[string]interface{}{"Table": "product", "Type": "*product.Entity", "Keyword": "remove", "Id": id}
			code, body := client.CallIt(i, "repoHandler")
			return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: code, Headers: r.Headers, Body: body})
		}

	case "find":
		if csv, ok := r.QueryStringParameters["ids"]; ok {

			ids := strings.Split(csv, ",")

			i := map[string]interface{}{"Table": "product", "Type": "*product.Entity", "Keyword": "find-many", "Ids": ids}
			code, body := client.CallIt(i, "repoHandler")
			if code != 200 {
				return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: code, Headers: r.Headers, Body: body})
			}

			var result []product.Entity
			_ = json.Unmarshal([]byte(body), &result)

			b, err := json.Marshal(result)
			if err != nil {
				return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: 200, Headers: r.Headers, Body: err.Error()})
			}

			return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: 200, Headers: r.Headers, Body: string(b)})

		} else if id, ok := r.QueryStringParameters["id"]; ok {

			i := map[string]interface{}{"Table": "product", "Type": "*product.Entity", "Keyword": "find-one", "Id": id}
			code, body := client.CallIt(i, "repoHandler")
			if code != 200 {
				return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: code, Headers: r.Headers, Body: body})
			}

			var result product.Entity
			_ = json.Unmarshal([]byte(body), &result)

			b, err := json.Marshal(result)
			if err != nil {
				return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: 200, Headers: r.Headers, Body: err.Error()})
			}

			return apigwp.HandleResponse(events.APIGatewayProxyResponse{StatusCode: 200, Headers: r.Headers, Body: string(b)})
		}

	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(Handle)
}
