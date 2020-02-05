package apigwp

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

type GoodValidator struct{}
type BadValidator struct{}
type BadMarshal struct {
	Foo string `json:"foo"`
	Bar string `json:"foo"`
}

var err = fmt.Errorf("error")

func (v *GoodValidator) Validate() error {
	return nil
}

func (v *BadValidator) Validate() error {
	return err
}

func TestRequestBadBody(t *testing.T) {
	if _, err := Request(events.APIGatewayProxyRequest{Body: "}{"}, &GoodValidator{}); err == nil {
		t.Fail()
	}
}

func TestRequestBadValidate(t *testing.T) {
	if _, err := Request(events.APIGatewayProxyRequest{Body: "{}"}, &BadValidator{}); err == nil {
		t.Fail()
	}
}

func TestResponseByte(t *testing.T) {
	if _, err := Response(200, []byte(``)); err != nil {
		t.Fail()
	}
}

func TestResponseString(t *testing.T) {
	if _, err := Response(200, ``); err != nil {
		t.Fail()
	}
}

func TestResponseError(t *testing.T) {
	if _, err := Response(200, err); err != nil {
		t.Fail()
	}
}

func TestResponseMarshalError(t *testing.T) {
	foo := make(chan bool)
	if _, err := Response(200, foo); err != nil {
		t.Fail()
	}
}

func TestResponseMarshal(t *testing.T) {
	if _, err := Response(200, BadMarshal{}); err != nil {
		t.Fail()
	}
}
