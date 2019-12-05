package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Address struct {
	Id      string `json:"id"`
	Session string `json:"session,omitempty"`
	Street1 string `json:"street_1,omitempty"`
	Street2 string `json:"street_2,omitempty"`
	UnitNum string `json:"unit_num,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip5    string `json:"zip_5,omitempty"`
	Zip4    string `json:"zip_4,omitempty"`
}

func (a *Address) Unmarshal(s string) error {
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		return err
	} else if a.Street1 == "" {
		return fmt.Errorf("bad street [%s]", a.Street1)
	} else if a.City == "" {
		return fmt.Errorf("bad city [%s]", a.City)
	} else if a.State == "" {
		return fmt.Errorf("bad state [%s]", a.State)
	} else if a.Zip5 == "" {
		return fmt.Errorf("bad zip [%s]", a.Zip5)
	} else if a.Id != "" {
		return nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		a.Id = id.String()
		return nil
	}
}
