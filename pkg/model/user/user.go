package user

import "os"

type Request struct {
	Op  string   `json:"op"`  // handler operation, add or delete
	Id  string   `json:"id"`  // id of user we are updating
	Col string   `json:"col"` // column name of ids to update
	Ids []string `json:"ids"` // ids to add or delete from Col
}

var Table = os.Getenv("USER_TABLE")
