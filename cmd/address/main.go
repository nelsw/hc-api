package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/client/faas/client"
	"hc-api/pkg/client/repo"
	"hc-api/pkg/factory/apigwp"
	"hc-api/pkg/model/address"
)

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var request address.Request

	ip, err := apigwp.Request(r, &request)
	if err != nil {
		return apigwp.Response(400, err)
	}

	request.SourceIp = ip // we set the ip here for CORS prevention

	if out, err := client.Call(request, "hcTokenHandler"); err != nil {
		return apigwp.Response(500, err)
	} else {
		_ = json.Unmarshal(out, &request)
		if len(request.JwtSlice) < 1 {
			return apigwp.Response(402)
		}
	}

	switch request.Op {

	case "find-one":
		if err := repo.FindOne(&request.Entity); err != nil {
			return apigwp.Response(404, err)
		}
		return apigwp.Response(200, &request.Entity)

	case "save":

		oldId := request.Id

		in := map[string]interface{}{"op": "validate", "address": request.Entity}

		if out, err := client.Call(&in, "hcUspsHandler"); err != nil {
			return apigwp.Response(500, err)
		} else {
			_ = json.Unmarshal(out, &request.Entity)
		}

		newId := base64.StdEncoding.EncodeToString([]byte(request.String()))

		if newId != oldId {

			ur1 := map[string]interface{}{"op": "add", "id": request.SourceId, "address_ids": []string{newId}}
			if _, err := client.Call(ur1, "hcUserHandler"); err != nil {
				return apigwp.Response(500, err)
			}

			if oldId != "" {

				ur2 := map[string]interface{}{"op": "delete", "id": request.SourceId, "address_ids": []string{oldId}}
				if _, err := client.Call(ur2, "hcUserHandler"); err != nil {
					return apigwp.Response(500, err)
				}
			}
		}

		return apigwp.Response(200, &request.Entity)
	}

	return apigwp.Response(400)
}

func main() {
	lambda.Start(Handle)
}
