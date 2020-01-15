package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"hc-api/internal/class"
	. "hc-api/service"
)

func Handle(r class.Request) ([]byte, error) {

	// Are we trying to find address entities by their id?
	if len(r.Ids) > 0 {
		if err := FindAllById(r.Table, r.Ids, &r.DataArr); err != nil {
			return nil, err
		} else {
			return json.Marshal(&r.DataArr)
		}
	}

	// We must be attempting to save an address entity.
	var a class.Address
	b, _ := json.Marshal(&r.Data)
	str, _ := VerifyAddress(string(b))
	_ = json.Unmarshal([]byte(str), &a)
	_ = Put(a, &r.Table)
	return json.Marshal(&r.Data)
}

func main() {
	lambda.Start(Handle)
}
