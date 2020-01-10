package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"os"
)

var client *lambda.Lambda

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	}); err != nil {
		log.Fatalf("Failed to connect to AWS: %s", err.Error())
	} else {
		client = lambda.New(sess)
	}
}

//func EntityRequest(name string, in, out interface{}) error {
//
//}

func ProxyRequest(name, ip string, in, out interface{}) error {
	b, _ := json.Marshal(in)
	payload, _ := json.Marshal(&events.APIGatewayProxyRequest{
		Body: string(b),
		RequestContext: events.APIGatewayProxyRequestContext{
			Identity: events.APIGatewayRequestIdentity{
				SourceIP: ip,
			},
		},
	})
	io, _ := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(name), Payload: payload})

	var resp events.APIGatewayProxyResponse
	_ = json.Unmarshal(io.Payload, &resp)
	return json.Unmarshal([]byte(resp.Body), out)
}

func invocation(name, body, ip string) (string, error) {
	var resp events.APIGatewayProxyResponse
	request := events.APIGatewayProxyRequest{
		Body: body,
		RequestContext: events.APIGatewayProxyRequestContext{
			Identity: events.APIGatewayRequestIdentity{
				SourceIP: ip,
			},
		},
	}
	if payload, err := json.Marshal(request); err != nil {
		return "", err
	} else if r, err := client.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(name),
		Payload:      payload,
	}); err != nil {
		return "", err
	} else if err := json.Unmarshal(r.Payload, &resp); err != nil {
		return "", err
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf(resp.Body)
	} else {
		return resp.Body, nil
	}
}
