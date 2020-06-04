package client

import (
	"encoding/json"
	"fmt"
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

func InvokeRaw(b []byte, s string) ([]byte, error) {
	if out, err := l.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(s),
		Payload:      b,
	}); err != nil {
		return nil, err
	} else {
		return out.Payload, nil
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

func Invoke(f string, i interface{}, o interface{}) error {
	r := events.APIGatewayProxyResponse{}
	if b, err := json.Marshal(&i); err != nil {
		return err
	} else if output, err := l.Invoke(&lambda.InvokeInput{FunctionName: aws.String(f), Payload: b}); err != nil {
		return err
	} else if err := json.Unmarshal(output.Payload, &r); err != nil {
		return err
	} else if r.StatusCode != 200 {
		return fmt.Errorf(r.Body)
	} else if o == nil {
		return nil
	} else {
		return json.Unmarshal([]byte(r.Body), &o)
	}
}
