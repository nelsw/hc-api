package client

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"os"
)

var l *lambda.Lambda

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}); err != nil {
		log.Fatalf("Failed to connect to AWS: %s", err.Error())
	} else {
		l = lambda.New(sess)
	}
}

func CallIt(i interface{}, s string) (int, string) {
	b, _ := json.Marshal(&i)
	input := lambda.InvokeInput{FunctionName: aws.String(s), Payload: b}
	if out, err := l.Invoke(&input); err != nil {
		return 500, err.Error()
	} else {
		return int(*out.StatusCode), string(out.Payload)
	}
}

func Invoke(f string, i interface{}) events.APIGatewayProxyResponse {
	r := events.APIGatewayProxyResponse{StatusCode: 500}
	b, _ := json.Marshal(&i)
	if output, err := l.Invoke(&lambda.InvokeInput{FunctionName: aws.String(f), Payload: b}); err != nil {
		r.Body = err.Error()
	} else if *output.StatusCode != 200 {
		r.Body = string(output.Payload)
	} else if err := json.Unmarshal(output.Payload, &r); err != nil {
		r.Body = err.Error()
	} else {
		_ = json.Unmarshal([]byte(r.Body), &r) // try it again jic
		r.StatusCode = 200
	}
	return r
}
