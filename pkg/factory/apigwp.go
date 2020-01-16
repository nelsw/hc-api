// Package represents a factory for creating API Gateway Proxy Request & Response values used in Proxy Lambda functions.
package factory

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/request"
)

// Required to transmit messages when CORS enabled in API Gateway.
var h = map[string]string{"Access-Control-Allow-Origin": "*"}

// Request attempts to unmarshal the wrapper body data into the given Validator.
// Returns nil unless an error occurs during interface deserialization or validation.
func Request(request events.APIGatewayProxyRequest, v request.Validator) error {
	body := request.Body
	fmt.Printf("REQUEST  body=[%s]\n", body)
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		return err
	} else if err := v.Validate(); err != nil {
		return err
	} else {
		return nil
	}
}

// Response returns an API Gateway Proxy Response with a nil error to provide detailed status codes and response bodies.
// While a status code must be provided, further arguments are recognized with reflection but not required.
func Response(i int, vv ...interface{}) (events.APIGatewayProxyResponse, error) {
	var body string
	if len(vv) > 0 {
		v := vv[0]
		if b, ok := v.([]byte); ok {
			body = string(b)
		} else if s, ok := v.(string); ok {
			body = s
		} else if err, ok := v.(error); ok {
			body = err.Error()
		} else if b, err := json.Marshal(v); err != nil {
			fmt.Printf("RESPONSE error, cannot marshal [%v]", v)
		} else {
			body = string(b)
		}
	}
	fmt.Printf("RESPONSE  code=[%d] body=[%s]\n", i, body)
	return events.APIGatewayProxyResponse{StatusCode: i, Headers: h, Body: body}, nil
}
