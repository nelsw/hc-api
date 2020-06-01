package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"sam-app/pkg/client/repo"
	"sam-app/pkg/factory/apigwp"
	"strings"
)

type Product struct {
	Id          string   `json:"id"`
	Category    string   `json:"category"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	ImageUrls   []string `json:"image_urls"`
	OwnerId     string   `json:"owner_id"`
	AddressId   string   `json:"address_id"` // shipping departure location
	Unit        string   `json:"unit"`       // LB, OZ, etc.
	Weight      int64    `json:"weight"`
	Stock       int8     `json:"stock"`
}

func (e *Product) ID() string {
	return e.Id
}

func (*Product) TableName() string {
	return table
}

func (e *Product) Validate() error {
	if len(e.Name) < 2 {
		return errName
	}
	return nil
}

var (
	table   = os.Getenv("PRODUCT_TABLE")
	errName = fmt.Errorf("product name must be at least 2 characters in length")
)

func Handle(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	apigwp.LogRequest(r)

	csv, ok := r.QueryStringParameters["ids"]
	if ok {
		ids := strings.Split(csv, ",")
		if out, err := repo.FindMany(&Product{}, ids); err != nil {
			return apigwp.Response(404, err)
		} else {
			return apigwp.Response(200, &out)
		}
	}

	return apigwp.Response(400, fmt.Errorf("...bad request, reached EOF.\n"))
}

func main() {
	lambda.Start(Handle)
}
