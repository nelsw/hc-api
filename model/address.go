package model

import (
	"fmt"
	"github.com/google/uuid"
)

type Address struct {
	Id      string `json:"id"`
	Street1 string `json:"street_1,omitempty"`
	Street2 string `json:"street_2,omitempty"`
	UnitNum string `json:"unit_num,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip5    string `json:"zip_5,omitempty"`
	Zip4    string `json:"zip_4,omitempty"`
}

func (address *Address) Validate() error {
	if address.Street1 == "" {
		return fmt.Errorf("bad street [%s]", address.Street1)
	} else if address.City == "" {
		return fmt.Errorf("bad city [%s]", address.City)
	} else if address.State == "" {
		return fmt.Errorf("bad state [%s]", address.State)
	} else if address.Zip5 == "" {
		return fmt.Errorf("bad zip [%s]", address.Zip5)
	} else {
		return nil
	}
}

func (address *Address) SetId() error {
	if address.Id != "" {
		return nil
	} else if id, err := uuid.NewUUID(); err != nil {
		return err
	} else {
		address.Id = id.String()
		return nil
	}
}
