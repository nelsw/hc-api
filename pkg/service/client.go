package service

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"hc-api/pkg/value"
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

// Finds by example - inspired by Hibernate ORM & Spring Result JPA.
func Find(in Invokable) error {
	return Request(in, "find")
}

func Save(in Invokable) error {
	return Request(in, "save")
}

func Update(in Invokable) error {
	return Request(in, "update")
}

func Request(in Invokable, s ...string) error {
	r := value.Request{Case: s[0], Table: *in.Table(), Ids: in.Ids(), Result: in}
	if b, err := json.Marshal(&r); err != nil {
		return err
	} else if _, err := l.Invoke(&lambda.InvokeInput{FunctionName: in.Function(), Payload: b}); err != nil {
		return err
	} else {
		return nil
	}
}

func Invoke(in Invokable) ([]byte, error) {
	if out, err := l.Invoke(&lambda.InvokeInput{FunctionName: in.Function(), Payload: in.Payload()}); err != nil {
		return nil, err
	} else {
		return out.Payload, nil
	}
}
