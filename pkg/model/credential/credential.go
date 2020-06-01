package credential

import (
	"os"
	"sam-app/pkg/util"
)

type Entity struct {
	Id         string `json:"id"`
	UserId     string `json:"user_id"`
	PasswordId string `json:"password_id"`
	Username   string `json:"username"` // transient
	Password   string `json:"password"` // transient
}

var table = os.Getenv("CREDENTIAL_TABLE")

func (e *Entity) Validate() error {
	if err := util.ValidateEmail(e.Id); err != nil {
		return err
	} else if err := util.ValidatePassword(e.Password); err != nil {
		return err
	}
	return nil
}

func (e *Entity) TableName() string {
	return table
}

func (e *Entity) ID() string {
	return e.Id
}
