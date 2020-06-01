package main

import (
	"github.com/aws/aws-lambda-go/events"
	"sam-app/test"
	"testing"
)

var (
	requestCtx = events.APIGatewayProxyRequestContext{
		Identity: events.APIGatewayRequestIdentity{
			SourceIP: test.Ip,
		},
	}
)

func TestHandleFindMany(t *testing.T) {
	ids := test.ProductId1 + "," + test.ProductId2
	r := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"ids": ids,
		},
		RequestContext: requestCtx,
	}

	out, _ := Handle(r)
	if out.StatusCode != 200 {
		t.Fail()
	}

}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
