package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"sam-app/pkg/model/product"
	"sam-app/test"
	"testing"
)

func TestHandleFindMany(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"path": "find",
			"ids":  test.ProductIds,
		},
	}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleFindOne(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"path": "find",
			"id":   test.ProductId,
		},
	}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleSaveOne(t *testing.T) {
	p := product.Entity{
		Id:          "10563954-a6b1-11ea-bbbf-26c506e4a61b",
		Category:    "Category Value",
		Name:        "Name Value",
		Description: "Description Value",
		Price:       99999999999,
		ImageUrls:   []string{"https://www.cbdrevolution.com/media/catalog/product/cache/3b283e46e55bcd65947f5adfccf62c98/c/r/cream_345.jpg"},
		OwnerId:     "OwnerId Value",
		AddressId:   "AddressId Value",
		Unit:        "Unit Value",
		Weight:      9999,
		Stock:       999999,
	}
	b, _ := json.Marshal(&p)
	body := base64.StdEncoding.EncodeToString(b)
	r := events.APIGatewayProxyRequest{
		Body: body,
		QueryStringParameters: map[string]string{
			"path": "save",
		},
		IsBase64Encoded: true,
	}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleDeleteOne(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"path": "remove",
			"id":   "fe7fab39-a46d-11ea-8817-2e51bfe26708",
		},
	}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleBadRequest(t *testing.T) {
	r := events.APIGatewayProxyRequest{}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
