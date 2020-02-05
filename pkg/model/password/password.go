package password

import (
	"golang.org/x/crypto/bcrypt"
	"os"
)

// Result structure for persisting and retrieving a users (encrypted) password. The user validation maintains a 1-1 FK
// relationship with the UserPassword validation for referential integrity, where user.PasswordId == userPassword.Value.
type Entity struct {
	Id   string `json:"id"`
	Hash string `json:"password"`
}

var table = os.Getenv("PASSWORD_TABLE")

func (*Entity) TableName() string {
	return table
}

func (e *Entity) ID() string {
	return e.Id
}

func (e *Entity) ComparePasswords(text string) error {
	return bcrypt.CompareHashAndPassword([]byte(e.Hash), []byte(text))
}
