package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/entity"
	"hc-api/pkg/service"
)

func Handle(e entity.Address) ([]byte, error) {

	switch e.Case {

	case "find":
		return service.Invoke(&e)

	case "save":
		e.Case = "verify"
		out, _ := service.Invoke(&e)
		_ = json.Unmarshal(out, &e)
		e.Case = "save"
		return service.Invoke(&e)
	}

	return nil, fmt.Errorf("bad case=[%s]", e.Case)
}

func main() {
	lambda.Start(Handle)
}
