package credential

import (
	"sam-app/pkg/util"
)

type Entity struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

func (e *Entity) Validate() error {
	if err := util.ValidateEmail(e.Id); err != nil {
		return err
	} else if err := util.ValidatePassword(e.Password); err != nil {
		return err
	} else {
		return nil
	}
}
