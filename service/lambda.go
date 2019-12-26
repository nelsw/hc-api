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
	CMD(string) Request
	Handler(string) Request
	Name(string) Request
	Body(interface{}) Request
	QSP(string, string) Request
	IP(string) Request
	Session(string) Request
	Build() (map[string]interface{}, error)
	Marshal(interface{}) error
}

type requestBuilder struct {
	name, ip string
	body     interface{}
	qsp      map[string]string
}

func (rb *requestBuilder) Handler(s string) Request {
	rb.name = "hc" + s + "Handler"
	return rb
}

func (rb *requestBuilder) Name(s string) Request {
	rb.name = s
	return rb
}

func (rb *requestBuilder) Body(i interface{}) Request {
	rb.body = i
	return rb
}

func (rb *requestBuilder) QSP(k, v string) Request {
	if rb.qsp == nil {
		rb.qsp = make(map[string]string)
	}
	rb.qsp[k] = v
	return rb
}

func (rb *requestBuilder) IP(s string) Request {
	rb.ip = s
	return rb
}

func (rb *requestBuilder) Session(s string) Request {
	rb.QSP("session", s)
	return rb
}

func (rb *requestBuilder) CMD(s string) Request {
	rb.QSP("cmd", s)
	return rb
}

func (rb *requestBuilder) Marshal(i interface{}) error {
	if b, err := json.Marshal(rb.body); err != nil {
		log.Println(rb.name)
		log.Println(rb.qsp)
		return err
	} else if s, err := invocation(rb.name, string(b), rb.ip, rb.qsp); err != nil {
		return err
	} else if err := json.Unmarshal([]byte(s), &i); err != nil {
		log.Println(s)
		log.Println(rb.name)
		log.Println(rb.qsp)
		panic(err)
		return err
	} else {
		return nil
	}
}

func (rb *requestBuilder) Build() (map[string]interface{}, error) {
	var body string
	if rb.body != nil {
		if b, err := json.Marshal(rb.body); err != nil {
			log.Println(rb.name)
			log.Println(rb.qsp)
			panic(err)
			return nil, err
		} else {
			body = string(b)
		}
	}
	if s, err := invocation(rb.name, body, rb.ip, rb.qsp); err != nil {
		log.Println(rb.name)
		log.Println(rb.qsp)
		panic(err)
		return nil, err
	} else if err := json.Unmarshal([]byte(s), &rb.body); err != nil {
		log.Println(s)
		log.Println(rb.name)
		log.Println(rb.qsp)
		panic(err)
		return nil, err
	} else {
		return rb.body.(map[string]interface{}), nil
	}
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
	return invocation(name, body, "", qsp)
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
	} else {
		return resp.Body, nil
	}
}
