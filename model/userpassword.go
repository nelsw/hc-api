package model

import "golang.org/x/crypto/bcrypt"

type UserPassword struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

func (up *UserPassword) Validate() error {
	if err := IsPasswordValid(up.Password); err != nil {
		return err
	} else if err := IsIdValid(up.Id); err != nil {
		return err
	} else {
		return nil
	}
}

func (up *UserPassword) PrePersist() error {
	if len(up.Password) > 24 {
		return nil
	} else if b, err := bcrypt.GenerateFromPassword([]byte(up.Password), bcrypt.MinCost); err != nil {
		return err
	} else {
		up.Password = string(b)
		return nil
	}
}
