package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"sam-app/pkg/client/repo/client"
	"sam-app/pkg/model/credential"
	"sam-app/pkg/model/product"
	"sam-app/pkg/model/profile"
	"sam-app/pkg/model/request"
	"sam-app/pkg/model/user"
	"strings"
)

var (
	InvalidKeyword = fmt.Errorf("bad keyword")
	InvalidType    = fmt.Errorf("bad type")
	typeRegistry   = map[string]interface{}{
		"*product.Entity":    product.Entity{},
		"*credential.Entity": credential.Entity{},
		"*user.Entity":       user.Entity{},
		"*profile.Entity":    profile.Entity{},
	}
)

func logRequest(i interface{}) {
	fmt.Printf("REQUEST value=[%v]\n", i)
}

func logResponse(i interface{}, e error) (interface{}, error) {
	fmt.Printf("RESPONSE value=[%v]\n", i)
	return i, e
}

func Handle(r request.Entity) (interface{}, error) {

	logRequest(r)

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
			return logResponse(nil, InvalidType)
		}

		switch r.Keyword {

		case "find-one":
			if err := client.FindById(r.Table, r.Id, &i); err != nil {
				return logResponse(nil, err)
			} else {
				return logResponse(i, nil)
			}

		case "find-many":
			if err := client.FindByIds(r.Table, r.Ids, &i); err != nil {
				return logResponse(nil, err)
			} else {
				return logResponse(i, nil)
			}

		case "find-all":
			if err := client.FindAll(r.Table, r.Attributes, &i); err != nil {
				return logResponse(nil, err)
			} else {
				return logResponse(i, nil)
			}
		}
	}

	return logResponse(nil, InvalidKeyword)
}

func main() {
	lambda.Start(Handle)
}
