package apigw

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

var headers = map[string]string{"Access-Control-Allow-Origin": "*"}

// Returns an APIGWPR with a 200 status code and a body determined by the type and value of i.
// If i is of type string, body == string value of i, else body == marshalled json value of i.
// See https://golang.org/ref/spec#Type_assertions for more information.
func Ok(i ...interface{}) (events.APIGatewayProxyResponse, error) {
	if len(i) > 0 {
		if s, ok := i[0].(string); ok {
			return response(200, s)
		}
		b, _ := json.Marshal(i[0]) // If i == nil, empty struct returned, see reflect.ValueOf(i interface{}).
		return response(200, string(b))
	}
	return response(200, "")
}

func BadRequest(i ...interface{}) (events.APIGatewayProxyResponse, error) {
	if len(i) > 0 {
		if s, ok := i[0].(string); ok {
			return response(400, s)
		} else if err, ok := i[0].(error); ok {
			return response(400, err.Error())
		}
	}
	return response(400, "")
}

func BadAuth(err error) (events.APIGatewayProxyResponse, error) {
	return response(401, err.Error())
}

func response(code int, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{StatusCode: code, Headers: headers, Body: body}, nil
}
