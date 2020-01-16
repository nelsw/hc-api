package entity

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"os"
)

// Result structure for persisting and retrieving a users (encrypted) password. The user entity maintains a 1-1 FK
// relationship with the UserPassword entity for referential integrity, where user.PasswordId == userPassword.Id.
type Password struct {
	Id      string `json:"id"`
	Encoded string `json:"password"`
	Decoded []byte `json:"password_dec"`
}

var passwordHandler = "hcPasswordHandler"
var passwordTable = os.Getenv("USER_PASSWORD_TABLE")

func (e *Password) Validate() error {
	return bcrypt.CompareHashAndPassword([]byte(e.Encoded), e.Decoded)
}

func (*Password) Function() *string {
	return &passwordHandler
}

func (e *Password) Payload() []byte {
	b, _ := json.Marshal(&e)
	return b
}

func (e *Password) Ids() []string {
	return []string{e.Id}
}

func (*Password) Table() *string {
	return &passwordTable
}
