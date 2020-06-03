package main

import (
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
	"sam-app/pkg/util"
	"strings"
	"sync"
)

var vs = []string{"UPS", "USPS", "FEDEX"}
var table = os.Getenv("TABLE")

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apigwp.LogRequest(r)
	var e order.Entity
	if err := json.Unmarshal([]byte(r.Body), &e); err != nil {
		return apigwp.Response(400, err)
	}

	tknRequest := events.APIGatewayProxyRequest{
		Path: "authenticate",
		QueryStringParameters: map[string]string{
			"token": r.QueryStringParameters["token"],
		},
	}

	if code, body := client.CallIt(tknRequest, "tokenHandler"); code != 200 {
		return apigwp.Response(code, body)
	}

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

		for i, p := range e.Packages {

			if p.Id == "" {
				continue
			}

			p.TotalLength = p.ProductLength
			p.TotalWeight = p.ProductWeight * float32(p.ProductQty)
			p.TotalHeight = p.ProductHeight * p.ProductQty
			p.TotalWidth = p.ProductWidth * p.ProductQty

			p.ZipOrigination = util.ZipFromAddressId(p.AddressId)
			p.ShipperStateCode = util.StateFromAddressId(p.AddressId)

			p.ZipDestination = util.ZipFromAddressId(e.AddressId)
			p.RecipientStateCode = util.StateFromAddressId(e.AddressId)

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

	case "save-order":
		// sum all packages
		for _, p := range e.Packages {
			e.OrderSum += p.ProductPrice + p.ShipRate
		}

		s, _ := uuid.NewUUID()
		e.Id = s.String()
		if err := repo.Save(table, e.Id, &e); err != nil {
			return apigwp.Response(500, err)
			//} else if err := repo.Update(&user, "add order_ids :p"); err != nil {
			//	return apigwp.Response(500, err)
		}

		userId := "" // need to grab from jwt
		ur1 := map[string]interface{}{"op": "add", "id": userId, "ids": []string{e.Id}, "keyword": "add order_ids"}
		if _, err := client.Call(ur1, "hcUserHandler"); err != nil {
			return apigwp.Response(500, err)
		}

		return apigwp.Response(200)
	}

	return apigwp.Response(400)
}

func main() {
	lambda.Start(HandleRequest)
}
