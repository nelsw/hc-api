package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/client/repo/client"
	"hc-api/pkg/model/address"
	"hc-api/pkg/model/offer"
	"hc-api/pkg/model/password"
	"hc-api/pkg/model/product"
	"hc-api/pkg/model/request"
	"strings"
)

var (
	InvalidKeyword = fmt.Errorf("bad keyword")
	InvalidType    = fmt.Errorf("bad type")
	typeRegistry   = map[string]interface{}{
		"*password.Entity": password.Entity{},
		"*product.Entity":  product.Entity{},
		"*address.Entity":  address.Entity{},
		"*offer.Entity":    offer.Entity{},
	}
)

func Handle(r request.Entity) (interface{}, error) {

	fmt.Printf("REQUEST  entity=[%v]\n", r)

	if r.Keyword == "save" {
		return nil, client.Save(r.Table, &r.Result)
	}

	if r.Keyword == "remove" {
		return nil, client.Remove(r.Table, r.Id, &r.Result)
	}

	if strings.Contains(r.Keyword, "add") || strings.Contains(r.Keyword, "delete") {
		return nil, client.Update(r.Table, r.Id, r.Keyword, r.Ids)
	}

	if strings.Contains(r.Keyword, "find-") {

		i, ok := typeRegistry[r.Type]
		if !ok {
			return nil, InvalidType
		}

		switch r.Keyword {

		case "find-one":
			if err := client.FindById(r.Table, r.Id, &i); err != nil {
				return nil, err
			} else {
				return i, nil
			}

		case "find-many":
			if err := client.FindByIds(r.Table, r.Ids, &i); err != nil {
				return nil, err
			} else {
				return i, nil
			}

		case "find-all":
			if err := client.FindAll(r.Table, r.Attributes, &i); err != nil {
				return nil, err
			} else {
				return i, nil
			}
		}
	}

	return nil, InvalidKeyword
}

func main() {
	lambda.Start(Handle)
}
