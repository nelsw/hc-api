package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/internal/class"
	"hc-api/internal/repository"
)

func HandleRequest(r class.Request) ([]byte, error) {

	switch r.Command {

	case "find-all":
		return repository.FindAll(r)

	case "find-by-brand-id":
		return repository.FindByBrandId(r)

	case "find-by-id":
		return repository.FindById(r)

	case "find-by-ids":
		return repository.FindByIds(r)

	case "save":
		return repository.Save(r)

	case "update":
		return repository.Update(r)
	}

	return nil, fmt.Errorf("bad command=[%s]", r.Command)
}

func main() {
	lambda.Start(HandleRequest)
}
