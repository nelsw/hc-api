package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"sam-app/pkg/model/product"
	"sam-app/test"
	"testing"
)

func TestHandleFindMany(t *testing.T) {
	r := events.APIGatewayProxyRequest{Path: "find", QueryStringParameters: map[string]string{"ids": test.ProductIds}}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleFindOne(t *testing.T) {
	r := events.APIGatewayProxyRequest{Path: "find", QueryStringParameters: map[string]string{"id": test.ProductId}}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleSaveOne(t *testing.T) {
	p := product.Entity{
		Category:    "Category Value",
		Name:        "Name Value",
		Description: "Description Value",
		Price:       99999999999,
		ImageUrls:   nil,
		OwnerId:     "OwnerId Value",
		AddressId:   "AddressId Value",
		Unit:        "Unit Value",
		Weight:      9999,
		Stock:       999999,
	}
	b, _ := json.Marshal(&p)
	r := events.APIGatewayProxyRequest{Path: "save", Body: string(b)}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleDeleteOne(t *testing.T) {
	r := events.APIGatewayProxyRequest{Path: "remove", QueryStringParameters: map[string]string{"id": test.ProductId}}
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
