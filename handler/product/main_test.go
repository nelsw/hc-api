package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"math/rand"
	"sam-app/pkg/model/product"
	"sam-app/test"
	"testing"
)

func TestHandleFindMany(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"path": "find",
			"ids":  test.ProductIds,
		},
	}); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleFindOne(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"path": "find",
			"id":   test.ProductId,
		},
	}); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleSaveOne(t *testing.T) {
	data, _ := json.Marshal(&product.Entity{
		"TestHandleDeleteOne",
		"OwnerId Value",
		"CBD Revolution",
		"Topical",
		"Soothe Skin Therapy Lotion",
		"Ideal for all skin types, perfectly balanced for use anytime. Especially effective after bathing and before engaging in outdoor activities. Deeply hydrates to alleviate itching.",
		"https://www.cbdrevolution.com/media/catalog/product/cache/3b283e46e55bcd65947f5adfccf62c98/s/k/skin_345.jpg",
		[]product.Option{
			{
				3495,
				170,
				"oz",
				rand.Intn(100),
				"",
				nil,
			},
		},
	})
	if out, _ := Handle(events.APIGatewayProxyRequest{
		Body: base64.StdEncoding.EncodeToString(data),
		QueryStringParameters: map[string]string{
			"path": "save",
		},
		IsBase64Encoded: true,
	}); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleDeleteOne(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"path": "remove",
			"id":   "TestHandleDeleteOne",
		},
	}); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleBadRequest(t *testing.T) {
	if out, _ := Handle(events.APIGatewayProxyRequest{}); out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
