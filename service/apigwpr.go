package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

var headers = map[string]string{"Access-Control-Allow-Origin": "*"}

type Headers map[string]string
type Object interface{}

type Response interface {
	Head(Headers) Response
	Code(int) Response
	Text(string) Response
	Data(Object) Response
	Error(error) Response
	Token(string) Response
	Build() (events.APIGatewayProxyResponse, error)
}

type builder struct {
	status  int
	headers Headers
	body    string
	token   string
	object  Object
}

func (b *builder) Head(headers Headers) Response {
	b.headers = headers
	return b
}

func (b *builder) Code(statusCode int) Response {
	b.status = statusCode
	return b
}

func (b *builder) Text(body string) Response {
	bytes, _ := json.Marshal(&map[string]string{"message": body})
	b.body = string(bytes)
	return b
}

func (b *builder) Data(object Object) Response {
	b.object = object
	return b
}

func (b *builder) Error(err error) Response {
	b.body = err.Error()
	return b
}

func (b *builder) Token(cookie string) Response {
	b.token = cookie
	return b
}

func (b *builder) Build() (events.APIGatewayProxyResponse, error) {
	if b.token != "" {
		b.headers["Set-Cookie"] = string(b.token)
	}
	if b.object != nil {
		// json package is resilient enough to marshal any non nil value, gl.
		bytes, _ := json.Marshal(b.object)
		b.body = string(bytes)
	}
	// we do this for clients that automatically unmarshal body as json, specifically jq.
	if b.body == "" {
		b.body = `{}`
	}
	r := events.APIGatewayProxyResponse{StatusCode: b.status, Headers: b.headers, Body: b.body}
	fmt.Printf("RESPONSE [%v]", r)
	return r, nil
}

func New() Response { return &builder{headers: headers} }

func Ok() Response { return &builder{status: http.StatusOK, headers: headers} }

func Unauthorized() Response {
	return &builder{status: http.StatusUnauthorized, headers: headers}
}

func BadRequest() Response {
	return &builder{status: http.StatusBadRequest, headers: headers}
}

func BadGateway() Response {
	return &builder{status: http.StatusBadGateway, headers: headers}
}

func InternalServerError() Response {
	return &builder{status: http.StatusInternalServerError, headers: headers}
}
