package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestHandle(t *testing.T) {
	out, _ := Handle(events.APIGatewayProxyRequest{})
	var got map[string]string
	_ = json.Unmarshal([]byte(out.Body), &got)
	for k, v := range m {
		if got[k] != v {
			t.Fail()
		}
	}
}

// for code coverage purposes only
func TestHandleMain(t *testing.T) {
	go main()
}
