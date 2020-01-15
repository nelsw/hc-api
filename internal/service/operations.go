package service

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/lambda"
	"hc-api/internal/class"
	"hc-api/internal/entity/base"
	"hc-api/internal/entity/token"
)

// Finds by example - inspired by Hibernate ORM & Spring Data JPA.
func Find(in base.Entity, ids ...string) error {
	if b, err := json.Marshal(&class.Request{Table: *in.Name(), Ids: ids, Data: in}); err != nil {
		return err
	} else if out, err := l.Invoke(&lambda.InvokeInput{FunctionName: in.Handler(), Payload: b}); err != nil {
		return err
	} else {
		return json.Unmarshal(out.Payload, &in)
	}
}

func Validate(in base.Entity) error {
	var resp error
	if out, err := l.Invoke(&lambda.InvokeInput{FunctionName: in.Handler(), Payload: in.Payload()}); err != nil {
		return err
	} else if err := json.Unmarshal(out.Payload, &resp); err != nil {
		return err
	} else {
		return resp
	}
}

func String(in base.Entity) (string, error) {
	if out, err := l.Invoke(&lambda.InvokeInput{FunctionName: in.Handler(), Payload: in.Payload()}); err != nil {
		return "", err
	} else {
		return string(out.Payload), nil
	}
}

func Auth(in token.Aggregate) (string, error) {
	if out, err := l.Invoke(&lambda.InvokeInput{FunctionName: in.Handler(), Payload: in.Payload()}); err != nil {
		return "", err
	} else {
		return string(out.Payload), nil
	}
}
