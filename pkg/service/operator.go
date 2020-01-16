package service

import "github.com/aws/aws-sdk-go/aws/request"

type Invokable interface {
	request.Validator
	Function() *string
	Payload() []byte
	Table() *string
	Ids() []string
}
