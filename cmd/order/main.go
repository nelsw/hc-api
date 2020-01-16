package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
	"hc-api/pkg/factory"
	"hc-api/pkg/service"
	"strings"
	"sync"
)

func HandleRequest(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	e := entity.Order{}
	if err := factory.Request(r, &e); err != nil {
		return factory.Response(400, err)
	}

	e.Authorization.SourceIp = r.RequestContext.Identity.SourceIP
	t := entity.Token{Authorization: e.Authorization}
	if uid, err := service.Invoke(&t); err != nil {
		return factory.Response(402, err)
	} else {
		e.UserId = string(uid)
	}

	switch e.Case {

	case "find-by-ids":
		if err := service.Find(&e); err != nil {
			return factory.Response(400, err)
		} else {
			return factory.Response(200, e.Results)
		}

	case "calc-rates":
		e.PackageIds = make([]string, len(e.Packages))
		for i, p := range e.Packages {
			if p.Id == "" {
				continue
			}
			e.PackageIds = append(e.PackageIds, p.Id)
			p.TotalLength = p.ProductLength
			p.TotalWeight = p.ProductWeight * float32(p.ProductQty)
			p.TotalHeight = p.ProductHeight * p.ProductQty
			p.TotalWidth = p.ProductWidth * p.ProductQty

			aFrom, _ := base64.StdEncoding.DecodeString(p.AddressIdFrom)
			arrFrom := strings.Split(string(aFrom), ", ")
			p.ZipOrigination = strings.Split(arrFrom[len(arrFrom)-2], "-")[0]
			p.ShipperStateCode = arrFrom[len(arrFrom)-3]

			aTo, _ := base64.StdEncoding.DecodeString(p.AddressIdTo)
			arrTo := strings.Split(string(aTo), ", ")
			p.ZipDestination = strings.Split(arrTo[len(arrTo)-2], "-")[0]
			p.RecipientStateCode = arrTo[len(arrTo)-3]

			e.Packages[i] = p

			vs := []string{"UPS", "USPS"}

			orders := make(chan entity.Order)

			var wg sync.WaitGroup
			wg.Add(len(vs))

			for _, v := range vs {
				go func(v string) {
					defer wg.Done()
					e.Vendor = v

					//service.Invoke()
					//err := internal.Invoke().Handler("Shipping").QSP("cmd", "rate").QSP("v", v).Body(o).Marshal(&o)
					//if err != nil {
					//	log.Println(err)
					//}
					orders <- e
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

			e.Rates = rates
			return factory.Response(200, &e)
		}

	case "save-order":
		// sum all packages
		for _, p := range e.Packages {
			e.OrderSum += p.ProductPrice + p.VendorPrice
		}

		user := entity.User{Id: e.UserId}
		if err := service.Find(&user); err != nil {
			return factory.Response(404, err)
		} else {
			e.ProfileId = user.ProfileId
		}

		profile := entity.Profile{Id: e.ProfileId}
		if err := service.Find(&profile); err != nil {
			return factory.Response(404, err)
		} else {
			e.Phone = profile.Phone
			e.Email = profile.Email
			e.FirstName = profile.FirstName
			e.LastName = profile.LastName
			e.PackageIds = nil
		}

		if err := service.Save(&e); err != nil {
			return factory.Response(500, err)
		} else if err := service.Update(&user, "add order_ids :p"); err != nil {
			return factory.Response(500, err)
		} else {
			return factory.Response(200, &e)
		}
	}
	return factory.Response(400, fmt.Sprintf("bad case=[%s]", e.Case))
}

func main() {
	lambda.Start(HandleRequest)
}
