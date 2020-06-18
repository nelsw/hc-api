// Package represents a factory for creating API Gateway Proxy BaseRequest & Response values used in Proxy Lambda functions.
package apigwp

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/request"
)

// Required to transmit messages when CORS enabled in API Gateway.
var defaultHeaders = map[string]string{"Access-Control-Allow-Origin": "*"}

const reqFmt = "request: {\n" +
	"\theaders: %v\n" +
	"\tpath: %s\n" +
	"\tquery_string_parameters: %v\n" +
	"\tbody: %s\n" +
	"}\n"

const resFmt = "response: {\n" +
	"\tcode: %d\n" +
	"\theaders: %v\n" +
	"\tbody: %s\n" +
	"}\n"

func body(v interface{}) string {
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
	fmt.Printf(reqFmt, r.Headers, r.Path, r.QueryStringParameters, r.Body)
}

func HandleResponse(r events.APIGatewayProxyResponse) (events.APIGatewayProxyResponse, error) {
	if r.Headers == nil {
		r.Headers = defaultHeaders
	} else {
		for k, v := range defaultHeaders {
			r.Headers[k] = v
		}
	}
	fmt.Printf(resFmt, r.StatusCode, r.Headers, r.Body)
	return r, nil
}

// BaseRequest attempts to unmarshal the wrapper body data into the given GoodValidator.
// Returns nil unless an error occurs during interface deserialization or validation.
func Request(r events.APIGatewayProxyRequest, i request.Validator) (string, error) {
	LogRequest(r)
	if err := json.Unmarshal([]byte(r.Body), &i); err != nil {
		return "", err
	} else if err := i.Validate(); err != nil {
		return "", err
	} else {
		return r.RequestContext.Identity.SourceIP, nil
	}
}

// Response returns an API Gateway Proxy Response with a nil error to provide detailed status codes and response bodies.
// While a status code must be provided, further arguments are recognized with reflection but not required.
func Response(i int, vv ...interface{}) (events.APIGatewayProxyResponse, error) {
	return ProxyResponse(i, defaultHeaders, vv)
}

func ProxyResponse(i int, h map[string]string, vv ...interface{}) (events.APIGatewayProxyResponse, error) {
	if len(vv) < 1 {
		return HandleResponse(events.APIGatewayProxyResponse{StatusCode: i, Headers: h})
	}
	body := body(vv[0])
	return HandleResponse(events.APIGatewayProxyResponse{StatusCode: i, Headers: h, Body: body})
}
