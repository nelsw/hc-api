package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/model/user"
)

var ErrBadOp = fmt.Errorf("bad operation\n")

func Handle(r user.Request) error {

	if r.Op == "add" || r.Op == "delete" {

		m := map[string]interface{}{
			"table":   r.TableName(),
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
