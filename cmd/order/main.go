package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/client/faas/client"
	"hc-api/pkg/client/repo"
	"hc-api/pkg/factory/apigwp"
	"hc-api/pkg/model/order"
	"hc-api/pkg/util"
	"strings"
	"sync"
)

var vs = []string{"UPS", "USPS", "FEDEX"}

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var request order.Request

	ip, err := apigwp.Request(r, &request)
	if err != nil {
		return apigwp.Response(400, err)
	}

	request.SourceIp = ip

	if out, err := client.Call(request, "hcTokenHandler"); err != nil {
		return apigwp.Response(500, err)
	} else {
		_ = json.Unmarshal(out, &request)
		if len(request.JwtSlice) < 1 {
			return apigwp.Response(402)
		}
	}

	switch request.Op {

	case "find-by-ids":
		if out, err := repo.FindMany(&request, request.Ids); err != nil {
			return apigwp.Response(400, err)
		} else {
			return apigwp.Response(200, &out)
		}

	case "calc-rates":
		//request.PackageIds = make([]string, len(request.Packages))

		for i, p := range request.Packages {

			if p.Id == "" {
				continue
			}

			p.TotalLength = p.ProductLength
			p.TotalWeight = p.ProductWeight * float32(p.ProductQty)
			p.TotalHeight = p.ProductHeight * p.ProductQty
			p.TotalWidth = p.ProductWidth * p.ProductQty

			p.ZipOrigination = util.ZipFromAddressId(p.AddressId)
			p.ShipperStateCode = util.StateFromAddressId(p.AddressId)

			p.ZipDestination = util.ZipFromAddressId(request.AddressId)
			p.RecipientStateCode = util.StateFromAddressId(request.AddressId)

			request.Packages[i] = p

			orders := make(chan order.Entity)

			var wg sync.WaitGroup
			wg.Add(len(vs))

			for _, v := range vs {

				go func(v string) {

					defer wg.Done()

					//proxy.Invoke()
					//err := pkg.Invoke().Handler("Shipping").QSP("cmd", "rate").QSP("v", v).Body(o).Marshal(&o)
					//if err != nil {
					//	log.Println(err)
					//}

					orders <- request
				}(v)
			}

			rates := map[string]map[string]map[string]string{}

			go func() {

				for o := range orders {

					for k, v := range o.Rates {

						if _, m := rates[k]; m {
							rates[k][o.Vendor] = v[o.Vendor]
						} else {
							rates[k] = v
						}
					}
				}
			}()

			wg.Wait()

			return apigwp.Response(200, &rates)
		}

	case "save-order":
		// sum all packages
		for _, p := range request.Packages {
			request.OrderSum += p.ProductPrice + p.ShipRate
		}

		if err := repo.SaveOne(&request.Entity); err != nil {
			return apigwp.Response(500, err)
			//} else if err := repo.Update(&user, "add order_ids :p"); err != nil {
			//	return apigwp.Response(500, err)
		} else {
			return apigwp.Response(200)
		}
	}

	return apigwp.Response(400)
}

func main() {
	lambda.Start(HandleRequest)
}
