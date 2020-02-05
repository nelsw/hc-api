package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/pkg/client/faas/client"
	"hc-api/pkg/model/user"
)

var ErrBadOp = fmt.Errorf("bad operation\n")

func Handle(r user.Request) error {

	if r.Op == "add" || r.Op == "delete" {

		m := map[string]interface{}{
			"table":   user.Table,
			"id":      r.Id,
			"ids":     r.Ids,
			"keyword": r.Op + " " + r.Col,
		}
		_, err := client.Call(&m, "hcRepoHandler")
		return err
	}

	return ErrBadOp
}

func main() {
	lambda.Start(Handle)
}
