package service

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
)

type LambdaRequest interface {
	CMD(string) LambdaRequest
	Handler(string) LambdaRequest
	Name(string) LambdaRequest
	Body(interface{}) LambdaRequest
	QSP(string, string) LambdaRequest
	Session(string) LambdaRequest
	IP(string) LambdaRequest
	ID(string) LambdaRequest
	Post() (string, error)
	Build() (map[string]interface{}, error)
	Marshal(interface{}) error
}

type requestBuilder struct {
	name, ip string
	body     interface{}
	qsp      map[string]string
}

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

func (rb *requestBuilder) Handler(s string) LambdaRequest {
	rb.name = "hc" + s + "Handler"
	return rb
}

func (rb *requestBuilder) Name(s string) LambdaRequest {
	rb.name = s
	return rb
}

func (rb *requestBuilder) Body(i interface{}) LambdaRequest {
	rb.body = i
	return rb
}

func (rb *requestBuilder) QSP(k, v string) LambdaRequest {
	if rb.qsp == nil {
		rb.qsp = make(map[string]string)
	}
	rb.qsp[k] = v
	return rb
}

func (rb *requestBuilder) IP(s string) LambdaRequest {
	rb.ip = s
	return rb
}

func (rb *requestBuilder) Session(s string) LambdaRequest {
	rb.QSP("session", s)
	return rb
}

func (rb *requestBuilder) CMD(s string) LambdaRequest {
	rb.QSP("cmd", s)
	return rb
}

func (rb *requestBuilder) ID(s string) LambdaRequest {
	rb.QSP("id", s)
	return rb
}

func (rb *requestBuilder) Marshal(i interface{}) error {
	if b, err := json.Marshal(&i); err != nil {
		panic(err)
	} else if s, err := invocation(rb.name, string(b), rb.ip, rb.qsp); err != nil {
		panic(err)
	} else if err := json.Unmarshal([]byte(s), &i); err != nil {
		panic(err)
	} else {
		return nil
	}
}

func (rb *requestBuilder) Build() (map[string]interface{}, error) {
	var body string
	if rb.body != nil {
		if b, err := json.Marshal(rb.body); err != nil {
			panic(err)
		} else {
			body = string(b)
		}
	}
	if s, err := invocation(rb.name, body, rb.ip, rb.qsp); err != nil {
		panic(err)
	} else if err := json.Unmarshal([]byte(s), &rb.body); err != nil {
		panic(err)
	} else {
		return rb.body.(map[string]interface{}), nil
	}
}

func (rb *requestBuilder) Post() (string, error) {
	if rb.body == nil {
		return invocation(rb.name, "", rb.ip, rb.qsp)
	} else if b, err := json.Marshal(rb.body); err != nil {
		panic(err)
	} else {
		return invocation(rb.name, string(b), rb.ip, rb.qsp)
	}
}

func Invoke() LambdaRequest { return &requestBuilder{} }

func VerifyCredentials(body string) (string, error) {
	return invocation("hcUsernameHandler", body, "", map[string]string{"cmd": "verify"})
}

func VerifyAddress(body string) (string, error) {
	return invocation("hcShippingHandler", body, "", map[string]string{"cmd": "verify"})
}

func NewSession(id, ip string) (string, error) {
	return invocation("hcSessionHandler", "", ip, map[string]string{"cmd": "create", "id": id, "ip": ip})
}

// deprecated, instead use Invoke().Handler("Session").Session(session).IP(ip).CMD("validate").Post()
func ValidateSession(sess, ip string) (string, error) {
	return invocation("hcSessionHandler", "", ip, map[string]string{"cmd": "validate", "session": sess, "ip": ip})
}

func invocation(name, body, ip string, qsp map[string]string) (string, error) {
	var resp events.APIGatewayProxyResponse
	request := events.APIGatewayProxyRequest{
		QueryStringParameters: qsp,
		Body:                  body,
		RequestContext: events.APIGatewayProxyRequestContext{
			Identity: events.APIGatewayRequestIdentity{
				SourceIP: ip,
			},
		},
	}
	if payload, err := json.Marshal(request); err != nil {
		return "", err
	} else if r, err := lc.Invoke(&lambda.InvokeInput{FunctionName: aws.String(name), Payload: payload}); err != nil {
		return "", err
	} else if err := json.Unmarshal(r.Payload, &resp); err != nil {
		return "", err
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf(resp.Body)
	} else {
		return resp.Body, nil
	}
}
