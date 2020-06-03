// This package is responsible for exporting generic request methods to domain specific request ƒ's.
package repo

import (
	"encoding/json"
	"sam-app/pkg/client/faas/client"
	"sam-app/pkg/model/request"
	"sam-app/pkg/util"
)

const functionName = "repoHandler"

type Storable interface {
	TableName() string
	ID() string
}

func Save(table, id string, i interface{}) error {
	return do(request.Entity{
		Id:         id,
		Type:       util.TypeOf(i),
		Table:      table,
		Ids:        nil,
		Attributes: nil,
		Keyword:    "save",
		Result:     i,
	}, i)
}

func SaveOne(i Storable) error {
	return do(request.Entity{
		Id:         i.ID(),
		Type:       util.TypeOf(i),
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
		Type:    util.TypeOf(i),
		Result:  i,
	}, i)
}

func FindById(table, id string, i interface{}) error {
	return do(request.Entity{
		Table:   table,
		Id:      id,
		Keyword: "find-one",
		Type:    util.TypeOf(i),
		Result:  i,
	}, i)
}

func FindByIds(table string, i interface{}, ids []string) ([]byte, error) {
	r := request.Entity{
		Table:   table,
		Ids:     ids,
		Keyword: "find-many",
		Type:    util.TypeOf(i),
	}
	payload, _ := json.Marshal(&r)
	return client.InvokeRaw(payload, functionName)
}

func FindMany(i Storable, ids []string) ([]byte, error) {
	r := request.Entity{
		Table:   i.TableName(),
		Ids:     ids,
		Keyword: "find-many",
		Type:    util.TypeOf(i),
	}
	payload, _ := json.Marshal(&r)
	return client.InvokeRaw(payload, functionName)
}

func Remove(i Storable, id string) error {
	r := request.Entity{
		Id:      id,
		Table:   i.TableName(),
		Keyword: "remove",
	}
	payload, _ := json.Marshal(&r)
	_, err := client.InvokeRaw(payload, functionName)
	return err
}

func Delete(i Storable, ids []string) error {
	r := request.Entity{
		Id:      i.ID(),
		Table:   i.TableName(),
		Ids:     ids,
		Keyword: "delete " + util.TypeName(i),
	}
	payload, _ := json.Marshal(&r)
	_, err := client.InvokeRaw(payload, functionName)
	return err
}

func do(r request.Entity, i interface{}) error {
	payload, _ := json.Marshal(&r)
	if b, err := client.InvokeRaw(payload, functionName); err != nil {
		return err
	} else {
		return json.Unmarshal(b, &i)
	}
}
