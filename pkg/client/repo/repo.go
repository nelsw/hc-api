// This package is responsible for exporting generic request methods to domain specific request ƒ's.
package repo

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"os"
	"reflect"
	"sam-app/pkg/model/request"
	"strings"
)

const functionName = "repoHandler"

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

type Storable interface {
	TableName() string
	ID() string
}

func Save(id string, i interface{}) error {
	return do(request.Entity{
		Id:         id,
		Type:       typeOf(i),
		Table:      typeName(i),
		Ids:        nil,
		Attributes: nil,
		Keyword:    "save",
		Result:     i,
	}, i)
}

func SaveOne(i Storable) error {
	return do(request.Entity{
		Id:         i.ID(),
		Type:       typeOf(i),
		Table:      i.TableName(),
		Ids:        nil,
		Attributes: nil,
		Keyword:    "save",
		Result:     i,
	}, i)
}

func FindOne(i Storable) error {
	return do(request.Entity{
		Table:   i.TableName(),
		Id:      i.ID(),
		Keyword: "find-one",
		Type:    typeOf(i),
	}, i)
}

func FindByIds(i interface{}, ids []string) ([]byte, error) {
	r := request.Entity{
		Table:   typeName(i),
		Ids:     ids,
		Keyword: "find-many",
		Type:    typeOf(i),
	}
	payload, _ := json.Marshal(&r)
	return InvokeRaw(payload, functionName)
}

func FindMany(i Storable, ids []string) ([]byte, error) {
	r := request.Entity{
		Table:   i.TableName(),
		Ids:     ids,
		Keyword: "find-many",
		Type:    typeOf(i),
	}
	payload, _ := json.Marshal(&r)
	return InvokeRaw(payload, functionName)
}

func Remove(i Storable, id string) error {
	r := request.Entity{
		Id:      id,
		Table:   i.TableName(),
		Keyword: "remove",
	}
	payload, _ := json.Marshal(&r)
	_, err := InvokeRaw(payload, functionName)
	return err
}

func do(r request.Entity, i interface{}) error {
	payload, _ := json.Marshal(&r)
	if b, err := InvokeRaw(payload, functionName); err != nil {
		return err
	} else {
		return json.Unmarshal(b, &i)
	}
}

func typeName(i interface{}) string {
	return strings.Split(typeOf(i), ".")[0]
}

func typeOf(v interface{}) string {
	return reflect.TypeOf(v).String()
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
