package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"os"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/client/repo"
	"sam-app/pkg/factory/apigwp"
	"sam-app/pkg/model/order"
	"strings"
	"sync"
)

var vs = []string{"UPS", "USPS", "FEDEX"}
var table = os.Getenv("TABLE")

func zipFromAddressId(s string) string {
	add, _ := base64.StdEncoding.DecodeString(s)
	csv := strings.Split(string(add), ", ")
	return strings.Split(csv[len(csv)-2], "-")[0]
}

func stateFromAddressId(s string) string {
	add, _ := base64.StdEncoding.DecodeString(s)
	csv := strings.Split(string(add), ", ")
	return csv[len(csv)-3]
}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	apigwp.LogRequest(r)

	authenticate := events.APIGatewayProxyRequest{Path: "authenticate", Headers: r.Headers}
	if authResponse := client.Invoke("tokenHandler", authenticate); authResponse.StatusCode != 200 {
		return apigwp.Response(401, authResponse.Body)
	} else {
		r.Headers = authResponse.Headers
	}

	e := order.Entity{}
	_ = json.Unmarshal([]byte(r.Body), &e)

	switch r.Path {

	case "find":
		if csv, ok := r.QueryStringParameters["ids"]; !ok {
			return apigwp.Response(400, "no ids")
		} else if out, err := repo.FindByIds(table, &e, strings.Split(csv, ",")); err != nil {
			return apigwp.Response(400, err)
		} else {
			return apigwp.Response(200, &out)
		}

	case "rates":

		if err := json.Unmarshal([]byte(r.Body), &e); err != nil {
			return apigwp.Response(400, err)
		}

		for i, p := range e.Packages {

			if p.Id == "" {
				continue
			}

			p.TotalLength = p.ProductLength
			p.TotalWeight = p.ProductWeight * float32(p.ProductQty)
			p.TotalHeight = p.ProductHeight * p.ProductQty
			p.TotalWidth = p.ProductWidth * p.ProductQty

			p.ZipOrigination = zipFromAddressId(p.AddressId)
			p.ShipperStateCode = stateFromAddressId(p.AddressId)

			p.ZipDestination = zipFromAddressId(e.AddressId)
			p.RecipientStateCode = stateFromAddressId(e.AddressId)

			e.Packages[i] = p

			ratesChan := make(chan order.Package)

			var wg sync.WaitGroup
			wg.Add(len(vs))

			for _, v := range vs {

				go func(v string) {

					defer wg.Done()

					//err := pkg.Invoke().Handler("Shipping").QSP("cmd", "rate").QSP("v", v).Body(o).Marshal(&o)
					//if err != nil {
					//	log.Println(err)
					//}

					ratesChan <- p
				}(v)
			}

			rates := map[string]map[string]map[string]interface{}{}

			go func() {

				for p := range ratesChan {

					if _, ok := rates[p.Id]; ok {
						rates[p.Id][p.ShipVendor] = map[string]interface{}{p.ShipService: p.ShipRate}
					} else {
						rates[p.Id] = map[string]map[string]interface{}{p.ShipVendor: {p.ShipService: p.ShipRate}}
					}
				}
			}()

			wg.Wait()

			return apigwp.Response(200, &rates)
		}

	case "save":

		for _, p := range e.Packages { // sum all packages
			e.OrderSum += p.ProductPrice + p.ShipRate
		}

		ids := uuid.New().String()
		if err := repo.Save(table, ids, &e); err != nil {
			return apigwp.Response(500, err)
		}

		add := events.APIGatewayProxyRequest{
			Path:                  "add",
			Headers:               r.Headers,
			QueryStringParameters: map[string]string{"id": e.UserId, "ids": ids, "col": "orders", "keyword": "add"},
		}
		if repoResponse := client.Invoke("userHandler", add); repoResponse.StatusCode != 200 {
			return apigwp.Response(500, repoResponse.Body)
		}

		return apigwp.Response(200)
	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(HandleRequest)
}
