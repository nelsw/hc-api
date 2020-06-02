package main

import (
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestHandleNoIds400(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		Path: "delete",
		QueryStringParameters: map[string]string{
			"col":   "orders",
			"token": "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBdWRpZW5jZSBWYWx1ZSIsImV4cCI6MTU5MTE3OTA1NywianRpIjoiSWQgVmFsdWUiLCJpYXQiOjE1OTEwOTI2NTcsImlzcyI6Iklzc3VlciBWYWx1ZSIsInN1YiI6IlN1YmplY3QgVmFsdWUifQ.SeLP6owuc0WPJqRMXZAUgorwsm2MhqC7_0C_-CPcMXU; Expires=Wed, 03 Jun 2020 10:10:57 GMT",
		},
	}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleNoCol400(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		Path: "delete",
		QueryStringParameters: map[string]string{
			"ids":   "foo,bar",
			"token": "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBdWRpZW5jZSBWYWx1ZSIsImV4cCI6MTU5MTE3OTA1NywianRpIjoiSWQgVmFsdWUiLCJpYXQiOjE1OTEwOTI2NTcsImlzcyI6Iklzc3VlciBWYWx1ZSIsInN1YiI6IlN1YmplY3QgVmFsdWUifQ.SeLP6owuc0WPJqRMXZAUgorwsm2MhqC7_0C_-CPcMXU; Expires=Wed, 03 Jun 2020 10:10:57 GMT",
		},
	}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleNoToken400(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		Path: "delete",
		QueryStringParameters: map[string]string{
			"ids": "foo,bar",
		},
	}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	}
}

func TestHandleAdd200(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		Path: "add",
		QueryStringParameters: map[string]string{
			"ids":   "foo,bar",
			"col":   "orders",
			"token": "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBdWRpZW5jZSBWYWx1ZSIsImV4cCI6MTU5MTE4MTI3MiwianRpIjoiMmY2YTdmM2YtYTQ2OS0xMWVhLWE2MWItNzY3ZTBlODViOTk1IiwiaWF0IjoxNTkxMDk0ODcyLCJpc3MiOiJJc3N1ZXIgVmFsdWUiLCJzdWIiOiJTdWJqZWN0IFZhbHVlIn0.wuFOf1mbNkgnPmz3_iIl-6UFlKw9AkO4IKkvDqFT4Tg; Expires=Wed, 03 Jun 2020 10:47:52 GMT",
		},
	}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleDelete200(t *testing.T) {
	r := events.APIGatewayProxyRequest{
		Path: "delete",
		QueryStringParameters: map[string]string{
			"ids":   "foo,bar",
			"col":   "orders",
			"token": "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBdWRpZW5jZSBWYWx1ZSIsImV4cCI6MTU5MTE4MTI3MiwianRpIjoiMmY2YTdmM2YtYTQ2OS0xMWVhLWE2MWItNzY3ZTBlODViOTk1IiwiaWF0IjoxNTkxMDk0ODcyLCJpc3MiOiJJc3N1ZXIgVmFsdWUiLCJzdWIiOiJTdWJqZWN0IFZhbHVlIn0.wuFOf1mbNkgnPmz3_iIl-6UFlKw9AkO4IKkvDqFT4Tg; Expires=Wed, 03 Jun 2020 10:47:52 GMT",
		},
	}
	if out, _ := Handle(r); out.StatusCode != 200 {
		t.Fail()
	}
}

func TestHandleNoPath400(t *testing.T) {
	r := events.APIGatewayProxyRequest{}
	if out, _ := Handle(r); out.StatusCode != 400 {
		t.Fail()
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
