package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"log"
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

	if token, ok := r.Headers["Authorize"]; !ok {
		return apigwp.Response(400, "missing token")
	} else {
		authenticate := events.APIGatewayProxyRequest{Path: "authenticate", Headers: map[string]string{"token": token}}
		if err := client.Invoke("tokenHandler", authenticate, &token); err != nil {
			return apigwp.Response(400, err)
		}
		r.Headers["Authorize"] = token
	}

	e := order.Entity{}

	switch r.Path {

	case "find-by-ids":
		if csv, ok := r.QueryStringParameters["ids"]; !ok {
			return apigwp.Response(400, fmt.Printf("bad qsp for ids [%s]", csv))
		} else if out, err := repo.FindByIds(table, e, strings.Split(csv, ",")); err != nil {
			return apigwp.Response(400, err)
		} else {
			return apigwp.Response(200, &out)
		}

	case "calc-rates":

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

					err := pkg.Invoke().Handler("Shipping").QSP("cmd", "rate").QSP("v", v).Body(o).Marshal(&o)
					if err != nil {
						log.Println(err)
					}

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

	case "save-order":

		id, ok := r.QueryStringParameters["id"]
		if !ok {
			return apigwp.Response(400, "no id provided")
		}

		for _, p := range e.Packages { // sum all packages
			e.OrderSum += p.ProductPrice + p.ShipRate
		}

		s, _ := uuid.NewUUID()
		e.Id = s.String()
		if err := repo.Save(table, e.Id, &e); err != nil {
			return apigwp.Response(500, err)
		} else if err := repo.Update(&user, "add order_ids :p"); err != nil {
			return apigwp.Response(500, err)
		}

		ur1 := map[string]interface{}{"op": "add", "id": id, "ids": []string{e.Id}, "keyword": "add order_ids"}
		if _, err := client.Call(ur1, "hcUserHandler"); err != nil {
			return apigwp.Response(500, err)
		}

		return apigwp.Response(200)
	}

	return apigwp.Response(400, fmt.Errorf("nothing returned for [%v].\n", r))
}

func main() {
	lambda.Start(HandleRequest)
}
