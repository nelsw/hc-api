// Package represents a factory for creating API Gateway Proxy BaseRequest & Response values used in Proxy Lambda functions.
package apigwp

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/request"
	"os"
)

// Required to transmit messages when CORS enabled in API Gateway.
var h = map[string]string{"Access-Control-Allow-Origin": "*"}
var stage = os.Getenv("STAGE")

func body(vv ...interface{}) string {
	if len(vv) < 1 {
		return ""
	}
	v := vv[0]
	if b, ok := v.([]byte); ok {
		return string(b)
	} else if s, ok := v.(string); ok {
		return s
	} else if err, ok := v.(error); ok {
		return err.Error()
	} else if b, err := json.Marshal(v); err != nil {
		return fmt.Sprintf("RESPONSE error, cannot marshal [%v]", v)
	} else {
		return string(b)
	}
}

func LogRequest(r events.APIGatewayProxyRequest) {
	fmt.Printf("PROXY REQUEST [%v]\n", r)
}

// BaseRequest attempts to unmarshal the wrapper body data into the given GoodValidator.
// Returns nil unless an error occurs during interface deserialization or validation.
func Request(r events.APIGatewayProxyRequest, i request.Validator) (string, error) {
	LogRequest(r)
	if err := json.Unmarshal([]byte(r.Body), &i); err != nil {
		return "", err
	} else if err := i.Validate(); err != nil {
		return "", err
	} else if stage == "TEST" || stage == "DEV" {
		return "127.0.0.1", nil
	} else {
		return r.RequestContext.Identity.SourceIP, nil
	}
}

// Response returns an API Gateway Proxy Response with a nil error to provide detailed status codes and response bodies.
// While a status code must be provided, further arguments are recognized with reflection but not required.
func Response(i int, vv ...interface{}) (events.APIGatewayProxyResponse, error) {
	return ProxyResponse(i, h, vv)
}

func ProxyResponse(i int, headers map[string]string, vv ...interface{}) (events.APIGatewayProxyResponse, error) {
	for k, v := range h {
		headers[k] = v
	}
	body := body(vv)
	fmt.Printf("PROXY RESPONSE | CODE=[%d] HEADERS=[%v] BODY=[%s]\n", i, headers, body)
	return events.APIGatewayProxyResponse{StatusCode: i, Headers: headers, Body: body}, nil
}
