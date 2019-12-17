package service

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
)

type Request interface {
	Handler(string) Request
	Name(string) Request
	Body(string) Request
	QSP(string, string) Request
	Build() (string, error)
}

type requestBuilder struct {
	name string
	body string
	qsp  map[string]string
}

func (rb *requestBuilder) Handler(s string) Request {
	rb.name = "hc" + s + "Handler"
	return rb
}

func (rb *requestBuilder) Name(s string) Request {
	rb.name = s
	return rb
}

func (rb *requestBuilder) Body(s string) Request {
	rb.body = s
	return rb
}

func (rb *requestBuilder) QSP(k, v string) Request {
	if rb.qsp == nil {
		rb.qsp = map[string]string{}
	}
	rb.qsp[k] = v
	return rb
}

func (rb *requestBuilder) Build() (string, error) {
	return invoke(rb.name, rb.body, rb.qsp)
}

func Invoke() Request { return &requestBuilder{} }

var lc *lambda.Lambda

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}); err != nil {
		log.Fatalf("Failed to connect to AWS: %s", err.Error())
	} else {
		lc = lambda.New(sess)
	}
}

func VerifyCredentials(body string) (string, error) {
	return invoke("hcUsernameHandler", body, map[string]string{"cmd": "verify"})
}

func VerifyPassword(p, id string) error {
	_, err := invoke("hcPasswordHandler", "", map[string]string{"cmd": "verify", "p": p, "id": id})
	return err
}

func VerifyAddress(body string) (string, error) {
	return invoke("hcShippingHandler", body, map[string]string{"cmd": "verify"})
}

func NewSession(id, ip string) (string, error) {
	return invoke("hcSessionHandler", "", map[string]string{"cmd": "create", "id": id, "ip": ip})
}

func ValidateSession(sess, ip string) (string, error) {
	return invoke("hcSessionHandler", "", map[string]string{"cmd": "validate", "session": sess, "ip": ip})
}

func invoke(name, body string, qsp map[string]string) (string, error) {
	var resp events.APIGatewayProxyResponse
	if payload, err := json.Marshal(events.APIGatewayProxyRequest{QueryStringParameters: qsp, Body: body}); err != nil {
		return "", err
	} else if r, err := lc.Invoke(&lambda.InvokeInput{FunctionName: aws.String(name), Payload: payload}); err != nil {
		return "", err
	} else if err := json.Unmarshal(r.Payload, &resp); err != nil {
		return "", err
	} else {
		return resp.Body, nil
	}
}
